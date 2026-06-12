package services

import (
	"errors"
	"fantasy-league-api/models"
	"fantasy-league-api/repository"
)

type PlayerService struct {
	playerRepo *repository.PlayerRepository
	teamRepo   *repository.TeamRepository
}

func NewPlayerService(playerRepo *repository.PlayerRepository, teamRepo *repository.TeamRepository) *PlayerService {
	return &PlayerService{
		playerRepo: playerRepo,
		teamRepo:   teamRepo,
	}
}

func (s *PlayerService) CreatePlayer(req models.CreatePlayerRequest) (*models.Player, error) {
	return s.playerRepo.CreatePlayer(req)
}

func (s *PlayerService) GetPlayers(page, pageSize int, sport, position string) ([]models.Player, error) {
	offset := (page - 1) * pageSize
	return s.playerRepo.GetPlayers(pageSize, offset, sport, position)
}

func (s *PlayerService) GetPlayerByID(id int) (*models.Player, error) {
	player, err := s.playerRepo.GetPlayerByID(id)
	if err != nil {
		return nil, err
	}
	if player == nil {
		return nil, errors.New("player not found")
	}
	return player, nil
}

func (s *PlayerService) AddPlayerToTeam(teamID, userID int, req models.AddPlayerToTeamRequest) error {

	// check team exists and belongs to user
	team, err := s.teamRepo.GetTeamByID(teamID)
	if err != nil {
		return err
	}
	if team == nil {
		return errors.New("team not found")
	}
	if team.OwnerID != userID {
		return errors.New("you do not own this team")
	}

	// check player exists
	player, err := s.playerRepo.GetPlayerByID(req.PlayerID)
	if err != nil {
		return err
	}
	if player == nil {
		return errors.New("player not found")
	}

	// check captain limit — only one captain per team
	if req.IsCaptain {
		captainCount, err := s.playerRepo.GetCaptainCount(teamID)
		if err != nil {
			return err
		}
		if captainCount >= 1 {
			return errors.New("team already has a captain")
		}
	}

	return s.playerRepo.AddPlayerToTeam(teamID, req.PlayerID, req.IsCaptain)
}

func (s *PlayerService) RemovePlayerFromTeam(teamID, playerID, userID int) error {

	// check team exists and belongs to user
	team, err := s.teamRepo.GetTeamByID(teamID)
	if err != nil {
		return err
	}
	if team == nil {
		return errors.New("team not found")
	}
	if team.OwnerID != userID {
		return errors.New("you do not own this team")
	}

	return s.playerRepo.RemovePlayerFromTeam(teamID, playerID)
}

func (s *PlayerService) GetTeamPlayers(teamID int) ([]models.TeamPlayerResponse, error) {
	team, err := s.teamRepo.GetTeamByID(teamID)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, errors.New("team not found")
	}

	return s.playerRepo.GetTeamPlayers(teamID)
}