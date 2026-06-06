package models

import "time"

type Team struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	LeagueID    int       `json:"league_id"`
	OwnerID     int       `json:"owner_id"`
	TotalPoints int       `json:"total_points"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateTeamRequest struct {
	Name string `json:"name" binding:"required,min=3,max=50"`
}

type TeamResponse struct {
	Team
	OwnerUsername string `json:"owner_username"`
	PlayerCount   int    `json:"player_count"`
}