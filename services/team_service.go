package services

import (
	"errors"
	"fantasy-league-api/models"
	"fantasy-league-api/repository"
)

type TeamService struct {
	teamRepo   *repository.TeamRepository
	leagueRepo *repository.LeagueRepository
}

func NewTeamService(teamRepo *repository.TeamRepository, leagueRepo *repository.LeagueRepository) *TeamService {
	return &TeamService{
		teamRepo:   teamRepo,
		leagueRepo: leagueRepo,
	}
}

func (s *TeamService) CreateTeam(req models.CreateTeamRequest, leagueID, userID int) (*models.Team, error) {

	// check league exists
	league, err := s.leagueRepo.GetLeagueByID(leagueID)
	if err != nil {
		return nil, err
	}
	if league == nil {
		return nil, errors.New("league not found")
	}

	// check league is still in draft
	if league.Status != "draft" {
		return nil, errors.New("league is no longer accepting teams")
	}

	// check user is a member or owner of the league
	if league.OwnerID != userID {
		isMember, err := s.leagueRepo.IsMember(leagueID, userID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, errors.New("you must join the league before creating a team")
		}
	}

	// check user doesn't already have a team in this league
	existing, err := s.teamRepo.GetTeamByLeagueAndOwner(leagueID, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("you already have a team in this league")
	}

	return s.teamRepo.CreateTeam(req.Name, leagueID, userID)
}

func (s *TeamService) GetTeamsByLeague(leagueID int) ([]models.TeamResponse, error) {
	league, err := s.leagueRepo.GetLeagueByID(leagueID)
	if err != nil {
		return nil, err
	}
	if league == nil {
		return nil, errors.New("league not found")
	}

	return s.teamRepo.GetTeamsByLeague(leagueID)
}

func (s *TeamService) GetTeamByID(id, userID int) (*models.Team, error) {
	team, err := s.teamRepo.GetTeamByID(id)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, errors.New("team not found")
	}

	return team, nil
}