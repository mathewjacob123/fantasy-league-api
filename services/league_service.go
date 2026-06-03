package services

import (
	"errors"
	"fantasy-league-api/models"
	"fantasy-league-api/repository"
)

type LeagueService struct {
	leagueRepo *repository.LeagueRepository
}

func NewLeagueService(leagueRepo *repository.LeagueRepository) *LeagueService {
	return &LeagueService{leagueRepo: leagueRepo}
}

func (s *LeagueService) CreateLeague(req models.CreateLeagueRequest, ownerID int) (*models.League, error) {
	return s.leagueRepo.CreateLeague(req.Name, ownerID, req.MaxTeams)
}

func (s *LeagueService) GetLeagues(page, pageSize int, userID int) ([]models.LeagueResponse, error) {
	offset := (page - 1) * pageSize
	return s.leagueRepo.GetLeagues(pageSize, offset, userID)
}

func (s *LeagueService) GetLeagueByID(id int, userID int) (*models.LeagueResponse, error) {
	league, err := s.leagueRepo.GetLeagueByID(id)
	if err != nil {
		return nil, err
	}
	if league == nil {
		return nil, errors.New("league not found")
	}

	memberCount, err := s.leagueRepo.GetMemberCount(league.ID)
	if err != nil {
		return nil, err
	}

	isMember, err := s.leagueRepo.IsMember(league.ID, userID)
	if err != nil {
		return nil, err
	}

	return &models.LeagueResponse{
		League:      *league,
		MemberCount: memberCount,
		IsMember:    isMember,
	}, nil
}

func (s *LeagueService) JoinLeague(leagueID, userID int) error {
	league, err := s.leagueRepo.GetLeagueByID(leagueID)
	if err != nil {
		return err
	}
	if league == nil {
		return errors.New("league not found")
	}

	if league.Status != "draft" {
		return errors.New("league is no longer accepting members")
	}

	isMember, err := s.leagueRepo.IsMember(leagueID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return errors.New("already a member of this league")
	}

	memberCount, err := s.leagueRepo.GetMemberCount(leagueID)
	if err != nil {
		return err
	}
	if memberCount >= league.MaxTeams {
		return errors.New("league is full")
	}

	return s.leagueRepo.JoinLeague(leagueID, userID)
}