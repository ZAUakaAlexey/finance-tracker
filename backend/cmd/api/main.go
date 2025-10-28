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

	r.Use(middlewares.ErrorHandler())

	api := r.Group("/api")
	{

		auth := api.Group("/auth")
		{
			auth.POST("/signup", handlers.Signup)
			auth.POST("/login", handlers.Login)
		}

		users := api.Group("/users")
		users.Use(middlewares.Authenticate)
		{
			users.GET("", handlers.GetUsers)           // GET /api/users?page=1&per_page=10
			users.GET("/search", handlers.SearchUsers) // GET /api/users/search?q=john&page=1&per_page=10
			users.GET("/me", handlers.GetCurrentUser)  // GET /api/users/me
			users.GET("/:id", handlers.GetUser)        // GET /api/users/:id
			users.PUT("/:id", handlers.UpdateUser)     // PUT /api/users/:id
			users.DELETE("/:id", handlers.DeleteUser)  // DELETE /api/users/:id
		}
	}

	r.NoRoute(middlewares.NotFound())

	log.Printf("Server starting on port %s", cfg.Port)
	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	r.Run(":" + cfg.Port)
}
