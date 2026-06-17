package repository

import (
	"database/sql"
	"fantasy-league-api/models"
)

type TeamRepository struct {
	db *sql.DB
}

func NewTeamRepository(db *sql.DB) *TeamRepository {
	return &TeamRepository{db: db}
}

func (r *TeamRepository) CreateTeam(name string, leagueID, ownerID int) (*models.Team, error) {
	query := `
		INSERT INTO teams (name, league_id, owner_id)
		VALUES ($1, $2, $3)
		RETURNING id, name, league_id, owner_id, total_points, created_at
	`

	team := &models.Team{}

	err := r.db.QueryRow(query, name, leagueID, ownerID).Scan(
		&team.ID,
		&team.Name,
		&team.LeagueID,
		&team.OwnerID,
		&team.TotalPoints,
		&team.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return team, nil
}

func (r *TeamRepository) GetTeamsByLeague(leagueID int) ([]models.TeamResponse, error) {
	query := `
    SELECT 
        t.id, t.name, t.league_id, t.owner_id,
        t.total_points, t.created_at,
        u.username as owner_username,
        COUNT(tp.id) as player_count
    FROM teams t
    JOIN users u ON t.owner_id = u.id
    LEFT JOIN team_players tp ON t.id = tp.team_id
    WHERE t.league_id = $1
    GROUP BY t.id, u.username
    ORDER BY t.total_points DESC
`

	rows, err := r.db.Query(query, leagueID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teams := []models.TeamResponse{}

	for rows.Next() {
		var team models.TeamResponse
		err := rows.Scan(
			&team.ID,
			&team.Name,
			&team.LeagueID,
			&team.OwnerID,
			&team.TotalPoints,
			&team.CreatedAt,
			&team.OwnerUsername,
			&team.PlayerCount,
		)
		if err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	return teams, nil
}

func (r *TeamRepository) GetTeamByID(id int) (*models.Team, error) {
	query := `
		SELECT id, name, league_id, owner_id, total_points, created_at
		FROM teams
		WHERE id = $1
	`

	team := &models.Team{}

	err := r.db.QueryRow(query, id).Scan(
		&team.ID,
		&team.Name,
		&team.LeagueID,
		&team.OwnerID,
		&team.TotalPoints,
		&team.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return team, nil
}

func (r *TeamRepository) GetTeamByLeagueAndOwner(leagueID, ownerID int) (*models.Team, error) {
	query := `
		SELECT id, name, league_id, owner_id, total_points, created_at
		FROM teams
		WHERE league_id = $1 AND owner_id = $2
	`

	team := &models.Team{}

	err := r.db.QueryRow(query, leagueID, ownerID).Scan(
		&team.ID,
		&team.Name,
		&team.LeagueID,
		&team.OwnerID,
		&team.TotalPoints,
		&team.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return team, nil
}