CREATE TABLE IF NOT EXISTS teams (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    league_id   INTEGER NOT NULL REFERENCES leagues(id) ON DELETE CASCADE,
    owner_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    total_points INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(league_id, owner_id)
);