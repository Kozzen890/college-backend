package httpapi

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"backend/internal/controllers"
	"backend/internal/middleware"
)

func SetupRouter(engine *gin.Engine, database *gorm.DB) {
	// Initialize controllers
	participantController := controllers.NewParticipantController(database)
	authController := controllers.NewAuthController(database)

	// Middleware global: set DB ke context agar bisa diakses di AuthMiddleware
	engine.Use(func(c *gin.Context) {
		c.Set("db", database)
		c.Next()
	})

	api := engine.Group("/api")
	{
		// Auth endpoints (tidak perlu login)
		api.POST("/login", authController.Login)
		api.POST("/refresh", authController.RefreshToken)

		// Participants endpoints
		api.POST("/participants", participantController.CreateParticipant) // Tidak perlu login untuk register

		// Protected endpoints (perlu login)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware())
		{
			// Auth endpoints yang perlu login
			protected.POST("/logout", authController.Logout)

			// Participants protected endpoints
			protected.GET("/participants", participantController.GetAllParticipants)
			protected.GET("/participants/count", participantController.CountParticipant)
			protected.GET("/participants/:id", participantController.GetParticipant)
			protected.PUT("/participants/:id", participantController.UpdateParticipant)
			protected.DELETE("/participants/:id", participantController.DeleteParticipant)

			// User endpoints
			protected.GET("/users/:id", authController.GetUserById)
			protected.GET("/admin/profile", authController.GetProfile)
		}
	}
}
