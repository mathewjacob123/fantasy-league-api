package models

import "time"

type League struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	OwnerID   int       `json:"owner_id"`
	MaxTeams  int       `json:"max_teams"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateLeagueRequest struct {
	Name     string `json:"name" binding:"required"`
	MaxTeams int    `json:"max_teams" binding:"required,min=2,max=20"`
}

type LeagueResponse struct {
	League
	MemberCount int `json:"member_count"`
	IsMember    bool `json:"is_member"`
}