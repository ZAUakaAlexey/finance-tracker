package main

import (
	"log"

	"github.com/ZAUakaAlexey/backend_go/internal/config"
	"github.com/ZAUakaAlexey/backend_go/internal/database"
	"github.com/ZAUakaAlexey/backend_go/internal/handlers"
	"github.com/ZAUakaAlexey/backend_go/internal/middlewares"
	"github.com/ZAUakaAlexey/backend_go/internal/validators"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	err = database.Connect(cfg)

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := validators.RegisterValidators(); err != nil {
		log.Fatal("Failed to register validators:", err)
	}

	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/health", handlers.Health)

		auth := api.Group("/auth")
		{
			auth.POST("/signup", handlers.Signup)
			auth.POST("/login", handlers.Login)
		}

		users := api.Group("/users")
		users.Use(middlewares.Authenticate)
		{
			users.GET("/:id", handlers.GetUser)
			users.PATCH("/:id", handlers.UpdateUser)
			users.DELETE("/:id")
			users.GET("/me", handlers.GetCurrentUser)
		}
	}

	log.Printf("Server starting on port %s", cfg.Port)
	r.Run(":" + cfg.Port)
}
