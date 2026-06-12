package repository

import (
	"database/sql"
	"fantasy-league-api/models"
)

type MatchRepository struct {
	db *sql.DB
}

func NewMatchRepository(db *sql.DB) *MatchRepository {
	return &MatchRepository{db: db}
}

func (r *MatchRepository) CreateMatch(req models.CreateMatchRequest) (*models.Match, error) {
	query := `
		INSERT INTO matches (
			home_team,
			away_team,
			sport,
			match_date
		)
		VALUES ($1, $2, $3, $4)
		RETURNING
			id,
			home_team,
			away_team,
			sport,
			status,
			match_date,
			created_at
	`

	match := &models.Match{}

	err := r.db.QueryRow(
		query,
		req.HomeTeam,
		req.AwayTeam,
		req.Sport,
		req.MatchDate,
	).Scan(
		&match.ID,
		&match.HomeTeam,
		&match.AwayTeam,
		&match.Sport,
		&match.Status,
		&match.MatchDate,
		&match.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return match, nil
}

func (r *MatchRepository) GetMatches(limit, offset int) ([]models.Match, error) {
	query := `
		SELECT
			id,
			home_team,
			away_team,
			sport,
			status,
			match_date,
			created_at
		FROM matches
		ORDER BY match_date DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	matches := []models.Match{}

	for rows.Next() {
		var match models.Match

		err := rows.Scan(
			&match.ID,
			&match.HomeTeam,
			&match.AwayTeam,
			&match.Sport,
			&match.Status,
			&match.MatchDate,
			&match.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		matches = append(matches, match)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return matches, nil
}

func (r *MatchRepository) GetMatchByID(id int) (*models.Match, error) {
	query := `
		SELECT
			id,
			home_team,
			away_team,
			sport,
			status,
			match_date,
			created_at
		FROM matches
		WHERE id = $1
	`

	match := &models.Match{}

	err := r.db.QueryRow(query, id).Scan(
		&match.ID,
		&match.HomeTeam,
		&match.AwayTeam,
		&match.Sport,
		&match.Status,
		&match.MatchDate,
		&match.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return match, nil
}