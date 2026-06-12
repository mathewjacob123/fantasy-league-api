CREATE TABLE IF NOT EXISTS players (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(255) NOT NULL,
    sport       VARCHAR(50) NOT NULL CHECK (sport IN ('football', 'cricket')),
    position    VARCHAR(50) NOT NULL,
    team_name   VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS team_players (
    id          SERIAL PRIMARY KEY,
    team_id     INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    player_id   INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    is_captain  BOOLEAN NOT NULL DEFAULT false,
    acquired_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(team_id, player_id)
);