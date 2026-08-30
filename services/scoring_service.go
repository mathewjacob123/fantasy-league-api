package services

import (
	"database/sql"
	"errors"
	"fantasy-league-api/models"
	"fantasy-league-api/repository"
)

type ScoringService struct {
	matchRepo    *repository.MatchRepository
	leagueRepo   *repository.LeagueRepository
	cacheService *CacheService
	db           *sql.DB
}

func NewScoringService(
	matchRepo *repository.MatchRepository,
	leagueRepo *repository.LeagueRepository,
	cacheService *CacheService,
	db *sql.DB,
) *ScoringService {
	return &ScoringService{
		matchRepo:    matchRepo,
		leagueRepo:   leagueRepo,
		cacheService: cacheService,
		db:           db,
	}
}

// calculatePoints — pure function, no DB calls
// easy to unit test independently
func calculatePoints(sport string, stats models.PlayerStats) int {
	points := 0

	if sport == "football" {
		points += stats.Goals * 6
		points += stats.Assists * 3
		points -= stats.YellowCards * 1
		points -= stats.RedCards * 3
		if stats.MinutesPlayed >= 90 {
			points += 2
		}
	}

	if sport == "cricket" {
		points += (stats.Runs / 10) * 1
		points += stats.Wickets * 25
		points += stats.Catches * 8
	}

	return points
}

func (s *ScoringService) SubmitMatchStats(matchID int, req models.SubmitStatsRequest) error {

	// get match details
	match, err := s.matchRepo.GetMatchByID(matchID)
	if err != nil {
		return err
	}
	if match == nil {
		return errors.New("match not found")
	}
	if match.Status == "completed" {
		return errors.New("stats already submitted for this match")
	}

	// START TRANSACTION — everything below is atomic
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// if anything fails, rollback everything
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// track points per team across all players
	teamPoints := map[int]int{}

	for _, stats := range req.Stats {

		// insert player stats for this match
		err = s.matchRepo.InsertPlayerStats(tx, matchID, stats)
		if err != nil {
			return err
		}

		// calculate points this player earned
		points := calculatePoints(match.Sport, stats)

		if points == 0 {
			continue
		}

		// find all fantasy teams that have this player
		teamIDs, err := s.matchRepo.GetTeamsWithPlayer(stats.PlayerID)
		if err != nil {
			return err
		}

		// accumulate points for each team
		for _, teamID := range teamIDs {
			teamPoints[teamID] += points
		}
	}

	// award points to each affected team
	for teamID, points := range teamPoints {

		// check if scorecard already exists — idempotency check
		exists, err := s.matchRepo.ScorecardExists(tx, teamID, matchID)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("points already awarded for this match")
		}

		// insert scorecard record
		err = s.matchRepo.InsertScorecard(tx, teamID, matchID, points)
		if err != nil {
			return err
		}

		// update team's running total
		err = s.matchRepo.UpdateTeamPoints(tx, teamID, points)
		if err != nil {
			return err
		}
	}

	// mark match as completed
	err = s.matchRepo.UpdateMatchStatus(matchID, "completed")
	if err != nil {
		return err
	}

	// invalidate leaderboard cache for all affected leagues
affectedLeagues := map[int]bool{}
for teamID := range teamPoints {
	team, err := s.leagueRepo.GetLeagueIDByTeamID(teamID)
	if err == nil && team != 0 {
		affectedLeagues[team] = true
	}
}
for leagueID := range affectedLeagues {
	s.cacheService.InvalidateLeaderboard(leagueID)
}

return tx.Commit()
}

func (s *ScoringService) CreateMatch(req models.CreateMatchRequest) (*models.Match, error) {
	return s.matchRepo.CreateMatch(req)
}

func (s *ScoringService) GetMatchByID(id int) (*models.Match, error) {
	match, err := s.matchRepo.GetMatchByID(id)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, errors.New("match not found")
	}
	return match, nil
}

func (s *ScoringService) GetMatchScores(matchID int) ([]models.Scorecard, error) {
	return s.matchRepo.GetMatchScores(matchID)
}