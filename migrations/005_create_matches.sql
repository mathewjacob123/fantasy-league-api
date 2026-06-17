CREATE TABLE IF NOT EXISTS matches (
    id          SERIAL PRIMARY KEY,
    home_team   VARCHAR(255) NOT NULL,
    away_team   VARCHAR(255) NOT NULL,
    sport       VARCHAR(50) NOT NULL CHECK (sport IN ('football', 'cricket')),
    match_date  TIMESTAMP NOT NULL,
    status      VARCHAR(20) NOT NULL DEFAULT 'scheduled' 
                CHECK (status IN ('scheduled', 'live', 'completed')),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS player_match_stats (
    id              SERIAL PRIMARY KEY,
    player_id       INTEGER NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    match_id        INTEGER NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    minutes_played  INTEGER NOT NULL DEFAULT 0,
    goals           INTEGER NOT NULL DEFAULT 0,
    assists         INTEGER NOT NULL DEFAULT 0,
    yellow_cards    INTEGER NOT NULL DEFAULT 0,
    red_cards       INTEGER NOT NULL DEFAULT 0,
    runs            INTEGER NOT NULL DEFAULT 0,
    wickets         INTEGER NOT NULL DEFAULT 0,
    catches         INTEGER NOT NULL DEFAULT 0,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(player_id, match_id)
);

CREATE TABLE IF NOT EXISTS scorecards (
    id          SERIAL PRIMARY KEY,
    team_id     INTEGER NOT NULL REFERENCES teams(id) ON DELETE CASCADE,
    match_id    INTEGER NOT NULL REFERENCES matches(id) ON DELETE CASCADE,
    points      INTEGER NOT NULL DEFAULT 0,
    calculated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(team_id, match_id)
);