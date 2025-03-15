package api

import (
	"github.com/gin-gonic/gin"
)

// SetupRoutes configures the API routes
func SetupRoutes(router *gin.Engine, handler *Handler) {
	// Health check
	router.GET("/health", handler.Health)

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Public API endpoints (User-Facing)

		// Badges endpoints
		v1.GET("/badges", handler.GetBadges)
		v1.GET("/badges/active", handler.GetActiveBadges)
		v1.GET("/badges/:id", handler.GetBadge)

		// User badges endpoints
		v1.GET("/users/:id/badges", handler.GetUserBadges)

		// Event processing endpoint
		v1.POST("/events", handler.ProcessEvent)

		// Admin API endpoints
		admin := v1.Group("/admin")
		{
			// Event types management
			admin.POST("/event-types", handler.CreateEventType)
			admin.GET("/event-types", handler.GetEventTypes)
			admin.GET("/event-types/:id", handler.GetEventType)
			admin.PUT("/event-types/:id", handler.UpdateEventType)
			admin.DELETE("/event-types/:id", handler.DeleteEventType)

			// Badge management
			admin.POST("/badges", handler.CreateBadge)
			admin.GET("/badges/:id/criteria", handler.GetBadgeWithCriteria)
			admin.PUT("/badges/:id", handler.UpdateBadge)
			admin.DELETE("/badges/:id", handler.DeleteBadge)

			// Condition types management
			admin.POST("/condition-types", handler.CreateConditionType)
			admin.GET("/condition-types", handler.GetConditionTypes)
			admin.GET("/condition-types/:id", handler.GetConditionType)
			admin.PUT("/condition-types/:id", handler.UpdateConditionType)
			admin.DELETE("/condition-types/:id", handler.DeleteConditionType)
		}
	}
}
