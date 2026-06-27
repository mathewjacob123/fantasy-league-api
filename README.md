# Fantasy League API

A production-style fantasy sports REST API built with Go, Gin, PostgreSQL, and Redis.
Supports football and cricket leagues with real-time scoring, leaderboards, and JWT authentication.

## Features

- JWT authentication with role-based access control (user/admin)
- League management — create, join, and manage fantasy leagues
- Team management — build and manage fantasy rosters
- Player management — admin-controlled player database with sport/position filtering
- Match scoring engine — transaction-based point calculation with idempotency protection
- Redis-cached leaderboard — 5 minute TTL with automatic invalidation on score updates
- Swagger API documentation

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Language | Go 1.21+ |
| Framework | Gin |
| Database | PostgreSQL |
| Cache | Redis (Memurai on Windows) |
| Auth | JWT (golang-jwt) |
| Password Hashing | bcrypt |
| Docs | Swagger (swaggo) |

## Architecture

```
Request → Middleware → Handler → Service → Repository → PostgreSQL
                                    ↓
                                 Redis Cache
```

Four layers — each with one job:
- **Handlers** — HTTP request/response only
- **Services** — business logic and validation
- **Repositories** — all database queries
- **Middleware** — JWT verification and role checking

## Project Structure

```
fantasy-league-api/
├── config/          # Environment config loader
├── db/              # PostgreSQL and Redis connections
├── handlers/        # HTTP handlers (auth, leagues, teams, players, matches)
├── middleware/       # JWT auth and admin middleware
├── migrations/      # SQL migration files
├── models/          # Request, response, and entity structs
├── repository/      # Database query layer
├── services/        # Business logic layer
└── docs/            # Auto-generated Swagger docs
```

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- Redis (or Memurai on Windows)

### Installation

```bash
git clone https://github.com/yourusername/fantasy-league-api
cd fantasy-league-api
go mod download
```

### Environment Setup

Create a `.env` file in the project root:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=fantasy_league
JWT_SECRET=your_secret_key
PORT=8080
REDIS_HOST=localhost
REDIS_PORT=6379
```

### Database Setup

Create the database in PostgreSQL:

```sql
CREATE DATABASE fantasy_league;
```

Run migrations in order from the `migrations/` folder:

```bash
# Run each file in pgAdmin or psql
migrations/001_create_users.sql
migrations/002_create_leagues.sql
migrations/003_create_teams.sql
migrations/004_create_players.sql
migrations/005_create_matches.sql
```

### Run

```bash
go run main.go
```

Server starts at `http://localhost:8080`

## API Documentation

Interactive Swagger docs available at:
```
http://localhost:8080/swagger/index.html
```

## Key API Endpoints

### Auth
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /auth/register | Register a new user |
| POST | /auth/login | Login and get JWT token |

### Leagues
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/leagues | Create a league |
| GET | /api/leagues | List all leagues (paginated) |
| GET | /api/leagues/:id | Get league details |
| POST | /api/leagues/:id/join | Join a league |
| GET | /api/leagues/:id/leaderboard | Get leaderboard (Redis cached) |

### Teams
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/leagues/:id/teams | Create a team |
| GET | /api/leagues/:id/teams | List teams in league |
| GET | /api/teams/:id | Get team details |

### Players
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | /api/players | List players (filterable) |
| GET | /api/players/:id | Get player details |
| POST | /api/teams/:id/players | Add player to team |
| DELETE | /api/teams/:id/players/:playerId | Remove player from team |
| GET | /api/teams/:id/players | Get team roster |

### Admin Only
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/admin/players | Create a player |
| POST | /api/admin/matches | Create a match |
| POST | /api/admin/matches/:id/stats | Submit match stats and trigger scoring |

## Scoring System

### Football
| Action | Points |
|--------|--------|
| Goal | +6 |
| Assist | +3 |
| 90 minutes played | +2 |
| Yellow card | -1 |
| Red card | -3 |

### Cricket
| Action | Points |
|--------|--------|
| Per 10 runs | +1 |
| Wicket | +25 |
| Catch | +8 |

## Testing

```bash
go test ./... -v
```

13 unit tests covering all scoring rules for football, cricket, and edge cases.

## What I Built and Learned

- Designed a normalised PostgreSQL schema across 8 tables with foreign keys, constraints, and indexes
- Implemented JWT authentication with role-based middleware (user/admin separation)
- Built a transaction-based scoring engine that atomically updates scores across multiple tables
- Added idempotency protection to prevent double-scoring on duplicate submissions
- Solved an N+1 query problem in the leagues list endpoint — reduced 201 queries to 1 using SQL JOINs
- Implemented Redis caching for the leaderboard with automatic cache invalidation on score updates
- Wrote 13 unit tests for the scoring engine using Go's table-driven test pattern