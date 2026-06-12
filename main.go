package main

import (
	"fantasy-league-api/config"
	"fantasy-league-api/db"
	"fantasy-league-api/handlers"
	"fantasy-league-api/middleware"
	"fantasy-league-api/repository"
	"fantasy-league-api/services"
	"log"
	"github.com/gin-gonic/gin"
	
)

func main() {
	cfg := config.Load()
	database := db.Connect(cfg)
	defer database.Close()

	
userRepo := repository.NewUserRepository(database)

	
authService := services.NewAuthService(userRepo, cfg.JWTSecret)

	
authHandler := handlers.NewAuthHandler(authService)

	
leagueRepo := repository.NewLeagueRepository(database)


leagueService := services.NewLeagueService(leagueRepo)


leagueHandler := handlers.NewLeagueHandler(leagueService)


teamRepo := repository.NewTeamRepository(database)


teamService := services.NewTeamService(teamRepo, leagueRepo)


teamHandler := handlers.NewTeamHandler(teamService)

// repositories
playerRepo := repository.NewPlayerRepository(database)

// services
playerService := services.NewPlayerService(playerRepo, teamRepo)

// handlers
playerHandler := handlers.NewPlayerHandler(playerService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	admin := r.Group("/api/admin")
admin.Use(middleware.AuthMiddleware(cfg.JWTSecret))
admin.Use(middleware.AdminMiddleware())
{
    admin.POST("/players", playerHandler.CreatePlayer)
}

	
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		api.GET("/me", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"user_id":  c.GetFloat64("user_id"),
				"email":    c.MustGet("email"),
				"username": c.MustGet("username"),
			})
		})
	api.POST("/leagues", leagueHandler.CreateLeague)
	api.GET("/leagues", leagueHandler.GetLeagues)
	api.GET("/leagues/:id", leagueHandler.GetLeagueByID)
	api.POST("/leagues/:id/join", leagueHandler.JoinLeague)
  api.POST("/leagues/:id/teams", teamHandler.CreateTeam)
  api.GET("/leagues/:id/teams", teamHandler.GetTeamsByLeague)
  api.GET("/teams/:id", teamHandler.GetTeamByID)

	api.GET("/players",                              playerHandler.GetPlayers)
    api.GET("/players/:id",                          playerHandler.GetPlayerByID)
    api.POST("/teams/:id/players",                   playerHandler.AddPlayerToTeam)
    api.DELETE("/teams/:id/players/:playerId",        playerHandler.RemovePlayerFromTeam)
    api.GET("/teams/:id/players",                    playerHandler.GetTeamPlayers)
	}

	log.Println("Server running on port", cfg.Port)
	r.Run(":" + cfg.Port)
}