package models

import "time"

type Match struct {
	ID        int       `json:"id"`
	HomeTeam  string    `json:"home_team"`
	AwayTeam  string    `json:"away_team"`
	Sport     string    `json:"sport"`
	MatchDate time.Time `json:"match_date"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateMatchRequest struct {
	HomeTeam  string    `json:"home_team"  binding:"required"`
	AwayTeam  string    `json:"away_team"  binding:"required"`
	Sport     string    `json:"sport"      binding:"required,oneof=football cricket"`
	MatchDate time.Time `json:"match_date" binding:"required"`
}

type PlayerStats struct {
	PlayerID      int `json:"player_id"      binding:"required"`
	MinutesPlayed int `json:"minutes_played"`
	Goals         int `json:"goals"`
	Assists       int `json:"assists"`
	YellowCards   int `json:"yellow_cards"`
	RedCards      int `json:"red_cards"`
	Runs          int `json:"runs"`
	Wickets       int `json:"wickets"`
	Catches       int `json:"catches"`
}

type SubmitStatsRequest struct {
	Stats []PlayerStats `json:"stats" binding:"required,min=1"`
}

type Scorecard struct {
	ID           int       `json:"id"`
	TeamID       int       `json:"team_id"`
	MatchID      int       `json:"match_id"`
	Points       int       `json:"points"`
	CalculatedAt time.Time `json:"calculated_at"`
}