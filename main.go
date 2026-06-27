package main

import (
	"fantasy-league-api/config"
	"fantasy-league-api/db"
	"fantasy-league-api/handlers"
	"fantasy-league-api/middleware"
	"fantasy-league-api/repository"
	"fantasy-league-api/services"
	"log"
	_ "fantasy-league-api/docs"
swaggerFiles "github.com/swaggo/files"
ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)
// @title           Fantasy League API
// @version         1.0
// @description     A fantasy sports league management API built with Go, Gin, PostgreSQL and Redis
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and your JWT token
func main() {
	// config and connections
	cfg := config.Load()
	database := db.Connect(cfg)
	redisClient := db.ConnectRedis(cfg)
	defer database.Close()

	// repositories
	userRepo   := repository.NewUserRepository(database)
	leagueRepo := repository.NewLeagueRepository(database)
	teamRepo   := repository.NewTeamRepository(database)
	playerRepo := repository.NewPlayerRepository(database)
	matchRepo  := repository.NewMatchRepository(database)

	// services
	cacheService  := services.NewCacheService(redisClient)
	authService   := services.NewAuthService(userRepo, cfg.JWTSecret)
	leagueService := services.NewLeagueService(leagueRepo, cacheService)
	teamService   := services.NewTeamService(teamRepo, leagueRepo)
	playerService := services.NewPlayerService(playerRepo, teamRepo)
	scoringService := services.NewScoringService(matchRepo, leagueRepo, cacheService, database)

	// handlers
	authHandler   := handlers.NewAuthHandler(authService)
	leagueHandler := handlers.NewLeagueHandler(leagueService)
	teamHandler   := handlers.NewTeamHandler(teamService)
	playerHandler := handlers.NewPlayerHandler(playerService)
	matchHandler  := handlers.NewMatchHandler(scoringService)

	r := gin.Default()

	// health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// public routes
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// admin routes
	admin := r.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	admin.Use(middleware.AdminMiddleware())
	{
		admin.POST("/players",           playerHandler.CreatePlayer)
		admin.POST("/matches",           matchHandler.CreateMatch)
		admin.POST("/matches/:id/stats", matchHandler.SubmitStats)
	}

	// protected routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// user
		api.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"user_id":  c.GetFloat64("user_id"),
				"email":    c.MustGet("email"),
				"username": c.MustGet("username"),
			})
		})

		// leagues
		api.POST("/leagues",              leagueHandler.CreateLeague)
		api.GET("/leagues",               leagueHandler.GetLeagues)
		api.GET("/leagues/:id",           leagueHandler.GetLeagueByID)
		api.POST("/leagues/:id/join",     leagueHandler.JoinLeague)
		api.GET("/leagues/:id/leaderboard", leagueHandler.GetLeaderboard)

		// teams
		api.POST("/leagues/:id/teams", teamHandler.CreateTeam)
		api.GET("/leagues/:id/teams",  teamHandler.GetTeamsByLeague)
		api.GET("/teams/:id",          teamHandler.GetTeamByID)

		// players
		api.GET("/players",                       playerHandler.GetPlayers)
		api.GET("/players/:id",                   playerHandler.GetPlayerByID)
		api.POST("/teams/:id/players",            playerHandler.AddPlayerToTeam)
		api.DELETE("/teams/:id/players/:playerId", playerHandler.RemovePlayerFromTeam)
		api.GET("/teams/:id/players",             playerHandler.GetTeamPlayers)

		// matches
		api.GET("/matches/:id",        matchHandler.GetMatch)
		api.GET("/matches/:id/scores", matchHandler.GetMatchScores)
	}

	log.Println("Server running on port", cfg.Port)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":" + cfg.Port)
}