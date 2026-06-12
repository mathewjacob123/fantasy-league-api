CREATE TABLE IF NOT EXISTS matches (
    id SERIAL PRIMARY KEY,

    home_team VARCHAR(255) NOT NULL,
    away_team VARCHAR(255) NOT NULL,

    sport VARCHAR(50) NOT NULL
        CHECK (sport IN ('football', 'cricket')),

    status VARCHAR(20) NOT NULL DEFAULT 'scheduled'
        CHECK (status IN ('scheduled', 'live', 'completed')),

    match_date TIMESTAMP NOT NULL,

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);