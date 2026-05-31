package repository

import (
	"database/sql"
	"fantasy-league-api/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(email, username, password string) (*models.User, error) {
	query := `
		INSERT INTO users (email, username, password)
		VALUES ($1, $2, $3)
		RETURNING id, email, username, created_at
	`

	user := &models.User{}

	err := r.db.QueryRow(query, email, username, password).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, email, username, password, created_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}

	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return user, nil
}