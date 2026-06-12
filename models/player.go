package models

import "time"

type Player struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Sport     string    `json:"sport"`
	Position  string    `json:"position"`
	TeamName  string    `json:"team_name"`
	CreatedAt time.Time `json:"created_at"`
}

type CreatePlayerRequest struct {
	Name     string `json:"name"      binding:"required"`
	Sport    string `json:"sport"     binding:"required,oneof=football cricket"`
	Position string `json:"position"  binding:"required"`
	TeamName string `json:"team_name" binding:"required"`
}

type AddPlayerToTeamRequest struct {
	PlayerID  int  `json:"player_id"  binding:"required"`
	IsCaptain bool `json:"is_captain"`
}

type TeamPlayerResponse struct {
	Player
	IsCaptain  bool      `json:"is_captain"`
	AcquiredAt time.Time `json:"acquired_at"`
}