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

	// repositories
	userRepo := repository.NewUserRepository(database)

	// services
	authService := services.NewAuthService(userRepo, cfg.JWTSecret)

	// handlers
	authHandler := handlers.NewAuthHandler(authService)

	// repositories
leagueRepo := repository.NewLeagueRepository(database)

// services
leagueService := services.NewLeagueService(leagueRepo)

// handlers
leagueHandler := handlers.NewLeagueHandler(leagueService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// auth routes — no middleware
	auth := r.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
	}

	// protected routes — requires JWT
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
	}

	log.Println("Server running on port", cfg.Port)
	r.Run(":" + cfg.Port)
}