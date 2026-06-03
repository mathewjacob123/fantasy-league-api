CREATE TABLE IF NOT EXISTS leagues (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) UNIQUE NOT NULL,
    owner_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    max_teams   INTEGER NOT NULL DEFAULT 10,
    status      VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'active', 'completed')),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS league_members (
    id          SERIAL PRIMARY KEY,
    league_id   INTEGER NOT NULL REFERENCES leagues(id) ON DELETE CASCADE,
    user_id     INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    joined_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(league_id, user_id)
);