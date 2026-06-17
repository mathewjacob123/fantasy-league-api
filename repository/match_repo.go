package repository

import (
	"database/sql"
	"fantasy-league-api/models"
)

type MatchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) *MatchRepository {
	return &MatchRepository{db: db}
}

func (r *MatchRepository) CreateMatch(req models.CreateMatchRequest) (*models.Match, error) {
	query := `
		INSERT INTO matches (home_team, away_team, sport, match_date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, home_team, away_team, sport, match_date, status, created_at
	`

	match := &models.Match{}

	err := r.db.QueryRow(query, req.HomeTeam, req.AwayTeam, req.Sport, req.MatchDate).Scan(
		&match.ID,
		&match.HomeTeam,
		&match.AwayTeam,
		&match.Sport,
		&match.MatchDate,
		&match.Status,
		&match.CreatedAt,
	)

	return match, err
}

func (r *MatchRepository) GetMatchByID(id int) (*models.Match, error) {
	query := `
		SELECT id, home_team, away_team, sport, match_date, status, created_at
		FROM matches WHERE id = $1
	`

	match := &models.Match{}

	err := r.db.QueryRow(query, id).Scan(
		&match.ID,
		&match.HomeTeam,
		&match.AwayTeam,
		&match.Sport,
		&match.MatchDate,
		&match.Status,
		&match.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return match, nil
}

func (r *MatchRepository) UpdateMatchStatus(id int, status string) error {
	query := `UPDATE matches SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *MatchRepository) InsertPlayerStats(tx *sql.Tx, matchID int, stats models.PlayerStats) error {
	query := `
		INSERT INTO player_match_stats 
		(player_id, match_id, minutes_played, goals, assists, yellow_cards, red_cards, runs, wickets, catches)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := tx.Exec(query,
		stats.PlayerID, matchID,
		stats.MinutesPlayed, stats.Goals, stats.Assists,
		stats.YellowCards, stats.RedCards,
		stats.Runs, stats.Wickets, stats.Catches,
	)

	return err
}

func (r *MatchRepository) GetTeamsWithPlayer(playerID int) ([]int, error) {
	query := `
		SELECT DISTINCT t.id
		FROM teams t
		JOIN team_players tp ON t.id = tp.team_id
		WHERE tp.player_id = $1
	`

	rows, err := r.db.Query(query, playerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teamIDs := []int{}

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		teamIDs = append(teamIDs, id)
	}

	return teamIDs, nil
}

func (r *MatchRepository) ScorecardExists(tx *sql.Tx, teamID, matchID int) (bool, error) {
	query := `SELECT COUNT(*) FROM scorecards WHERE team_id = $1 AND match_id = $2`
	var count int
	err := tx.QueryRow(query, teamID, matchID).Scan(&count)
	return count > 0, err
}

func (r *MatchRepository) InsertScorecard(tx *sql.Tx, teamID, matchID, points int) error {
	query := `
		INSERT INTO scorecards (team_id, match_id, points)
		VALUES ($1, $2, $3)
	`
	_, err := tx.Exec(query, teamID, matchID, points)
	return err
}

func (r *MatchRepository) UpdateTeamPoints(tx *sql.Tx, teamID, points int) error {
	query := `UPDATE teams SET total_points = total_points + $1 WHERE id = $2`
	_, err := tx.Exec(query, points, teamID)
	return err
}

func (r *MatchRepository) GetMatchScores(matchID int) ([]models.Scorecard, error) {
	query := `
		SELECT id, team_id, match_id, points, calculated_at
		FROM scorecards
		WHERE match_id = $1
		ORDER BY points DESC
	`

	rows, err := r.db.Query(query, matchID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scorecards := []models.Scorecard{}

	for rows.Next() {
		var s models.Scorecard
		err := rows.Scan(&s.ID, &s.TeamID, &s.MatchID, &s.Points, &s.CalculatedAt)
		if err != nil {
			return nil, err
		}
		scorecards = append(scorecards, s)
	}

	return scorecards, nil
}