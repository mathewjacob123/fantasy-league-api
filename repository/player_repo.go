package repository

import (
	"database/sql"
	"fantasy-league-api/models"
	"fmt"
)

type PlayerRepository struct {
	db *sql.DB
}

func NewPlayerRepository(db *sql.DB) *PlayerRepository {
	return &PlayerRepository{db: db}
}

func (r *PlayerRepository) CreatePlayer(req models.CreatePlayerRequest) (*models.Player, error) {
	query := `
		INSERT INTO players (name, sport, position, team_name)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, sport, position, team_name, created_at
	`

	player := &models.Player{}

	err := r.db.QueryRow(query, req.Name, req.Sport, req.Position, req.TeamName).Scan(
		&player.ID,
		&player.Name,
		&player.Sport,
		&player.Position,
		&player.TeamName,
		&player.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return player, nil
}

func (r *PlayerRepository) GetPlayers(limit, offset int, sport, position string) ([]models.Player, error) {
	baseQuery := `
		SELECT id, name, sport, position, team_name, created_at
		FROM players
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if sport != "" {
		baseQuery += fmt.Sprintf(" AND sport = $%d", argCount)
		args = append(args, sport)
		argCount++
	}

	if position != "" {
		baseQuery += fmt.Sprintf(" AND position = $%d", argCount)
		args = append(args, position)
		argCount++
	}

	baseQuery += fmt.Sprintf(" ORDER BY name ASC LIMIT $%d OFFSET $%d", argCount, argCount+1)
	args = append(args, limit, offset)

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := []models.Player{}

	for rows.Next() {
		var player models.Player
		err := rows.Scan(
			&player.ID,
			&player.Name,
			&player.Sport,
			&player.Position,
			&player.TeamName,
			&player.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		players = append(players, player)
	}

	return players, nil
}

func (r *PlayerRepository) GetPlayerByID(id int) (*models.Player, error) {
	query := `
		SELECT id, name, sport, position, team_name, created_at
		FROM players WHERE id = $1
	`

	player := &models.Player{}

	err := r.db.QueryRow(query, id).Scan(
		&player.ID,
		&player.Name,
		&player.Sport,
		&player.Position,
		&player.TeamName,
		&player.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return player, nil
}

func (r *PlayerRepository) AddPlayerToTeam(teamID, playerID int, isCaptain bool) error {
	query := `
		INSERT INTO team_players (team_id, player_id, is_captain)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.Exec(query, teamID, playerID, isCaptain)
	return err
}

func (r *PlayerRepository) RemovePlayerFromTeam(teamID, playerID int) error {
	query := `
		DELETE FROM team_players
		WHERE team_id = $1 AND player_id = $2
	`

	result, err := r.db.Exec(query, teamID, playerID)
	if err != nil {
		return err
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("player not found in team")
	}

	return nil
}

func (r *PlayerRepository) GetTeamPlayers(teamID int) ([]models.TeamPlayerResponse, error) {
	query := `
		SELECT 
			p.id, p.name, p.sport, p.position, p.team_name, p.created_at,
			tp.is_captain, tp.acquired_at
		FROM team_players tp
		JOIN players p ON tp.player_id = p.id
		WHERE tp.team_id = $1
		ORDER BY tp.is_captain DESC, p.name ASC
	`

	rows, err := r.db.Query(query, teamID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	players := []models.TeamPlayerResponse{}

	for rows.Next() {
		var p models.TeamPlayerResponse
		err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Sport,
			&p.Position,
			&p.TeamName,
			&p.CreatedAt,
			&p.IsCaptain,
			&p.AcquiredAt,
		)
		if err != nil {
			return nil, err
		}
		players = append(players, p)
	}

	return players, nil
}

func (r *PlayerRepository) GetCaptainCount(teamID int) (int, error) {
	query := `
		SELECT COUNT(*) FROM team_players
		WHERE team_id = $1 AND is_captain = true
	`

	var count int
	err := r.db.QueryRow(query, teamID).Scan(&count)
	return count, err
}