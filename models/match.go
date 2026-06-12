package models

import "time"

type Match struct {
	ID        int       `json:"id"`
	HomeTeam  string    `json:"home_team"`
	AwayTeam  string    `json:"away_team"`
	Sport     string    `json:"sport"`
	Status    string    `json:"status"`
	MatchDate time.Time `json:"match_date"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateMatchRequest struct {
	HomeTeam string    `json:"home_team" binding:"required"`
	AwayTeam string    `json:"away_team" binding:"required"`
	Sport    string    `json:"sport" binding:"required,oneof=football cricket"`
	MatchDate time.Time `json:"match_date" binding:"required"`
}