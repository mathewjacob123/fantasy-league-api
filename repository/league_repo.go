package repository

import (
	"database/sql"
	"fantasy-league-api/models"
)

type LeagueRepository struct {
	db *sql.DB
}

func NewLeagueRepository(db *sql.DB) *LeagueRepository {
	return &LeagueRepository{db: db}
}

func (r *LeagueRepository) CreateLeague(name string, ownerID int, maxTeams int) (*models.League, error) {
	query := `
		INSERT INTO leagues (name, owner_id, max_teams)
		VALUES ($1, $2, $3)
		RETURNING id, name, owner_id, max_teams, status, created_at
	`

	league := &models.League{}

	err := r.db.QueryRow(query, name, ownerID, maxTeams).Scan(
		&league.ID,
		&league.Name,
		&league.OwnerID,
		&league.MaxTeams,
		&league.Status,
		&league.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return league, nil
}

func (r *LeagueRepository) GetLeagues(limit, offset, userID int) ([]models.LeagueResponse, error) {
	query := `
		SELECT 
			l.id,
			l.name,
			l.owner_id,
			l.max_teams,
			l.status,
			l.created_at,
			COUNT(lm.id) as member_count,
			MAX(CASE WHEN lm.user_id = $3 THEN 1 ELSE 0 END) as is_member
		FROM leagues l
		LEFT JOIN league_members lm ON l.id = lm.league_id
		GROUP BY l.id
		ORDER BY l.created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leagues := []models.LeagueResponse{}

	for rows.Next() {
		var league models.LeagueResponse
		var isMember int

		err := rows.Scan(
			&league.ID,
			&league.Name,
			&league.OwnerID,
			&league.MaxTeams,
			&league.Status,
			&league.CreatedAt,
			&league.MemberCount,
			&isMember,
		)
		if err != nil {
			return nil, err
		}

		league.IsMember = isMember == 1
		leagues = append(leagues, league)
	}

	return leagues, nil
}

func (r *LeagueRepository) GetLeagueByID(id int) (*models.League, error) {
	query := `
		SELECT id, name, owner_id, max_teams, status, created_at
		FROM leagues
		WHERE id = $1
	`

	league := &models.League{}

	err := r.db.QueryRow(query, id).Scan(
		&league.ID,
		&league.Name,
		&league.OwnerID,
		&league.MaxTeams,
		&league.Status,
		&league.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return league, nil
}

func (r *LeagueRepository) JoinLeague(leagueID, userID int) error {
	query := `
		INSERT INTO league_members (league_id, user_id)
		VALUES ($1, $2)
	`

	_, err := r.db.Exec(query, leagueID, userID)
	return err
}

func (r *LeagueRepository) GetMemberCount(leagueID int) (int, error) {
	query := `
		SELECT COUNT(*) FROM league_members
		WHERE league_id = $1
	`

	var count int
	err := r.db.QueryRow(query, leagueID).Scan(&count)
	return count, err
}

func (r *LeagueRepository) IsMember(leagueID, userID int) (bool, error) {
	query := `
		SELECT COUNT(*) FROM league_members
		WHERE league_id = $1 AND user_id = $2
	`

	var count int
	err := r.db.QueryRow(query, leagueID, userID).Scan(&count)
	return count > 0, err
}

func (r *LeagueRepository) GetLeaderboard(leagueID int) ([]models.TeamResponse, error) {
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

func (r *LeagueRepository) GetLeagueIDByTeamID(teamID int) (int, error) {
	query := `SELECT league_id FROM teams WHERE id = $1`
	var leagueID int
	err := r.db.QueryRow(query, teamID).Scan(&leagueID)
	return leagueID, err
}