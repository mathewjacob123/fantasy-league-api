package services

import (
	"context"
	"encoding/json"
	"fantasy-league-api/models"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type CacheService struct {
	client *redis.Client
}

func NewCacheService(client *redis.Client) *CacheService {
	return &CacheService{client: client}
}

func leaderboardKey(leagueID int) string {
	return fmt.Sprintf("leaderboard:league:%d", leagueID)
}

func (s *CacheService) GetLeaderboard(leagueID int) ([]models.TeamResponse, error) {
	if s.client == nil {
		return nil, nil
	}

	ctx := context.Background()
	data, err := s.client.Get(ctx, leaderboardKey(leagueID)).Bytes()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var teams []models.TeamResponse
	err = json.Unmarshal(data, &teams)
	return teams, err
}

func (s *CacheService) SetLeaderboard(leagueID int, teams []models.TeamResponse) error {
	if s.client == nil {
		return nil
	}

	ctx := context.Background()
	data, err := json.Marshal(teams)
	if err != nil {
		return err
	}

	return s.client.Set(ctx, leaderboardKey(leagueID), data, 5*time.Minute).Err()
}

func (s *CacheService) InvalidateLeaderboard(leagueID int) error {
	if s.client == nil {
		return nil
	}

	ctx := context.Background()
	return s.client.Del(ctx, leaderboardKey(leagueID)).Err()
}