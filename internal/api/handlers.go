package api

import (
	"net/http"
	"strconv"

	"github.com/badge-assignment-system/internal/models"
	"github.com/badge-assignment-system/internal/service"
	"github.com/gin-gonic/gin"
)

// Handler handles HTTP requests
type Handler struct {
	Service *service.Service
}

// NewHandler creates a new handler
func NewHandler(s *service.Service) *Handler {
	return &Handler{
		Service: s,
	}
}

// respondWithError responds with a JSON error message
func respondWithError(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"error": message,
	})
}

// Health checks the health of the API
func (h *Handler) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "healthy",
	})
}

// CreateEventType handles creating a new event type
func (h *Handler) CreateEventType(c *gin.Context) {
	var req models.NewEventTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	eventType, err := h.Service.CreateEventType(&req)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, eventType)
}

// GetEventTypes handles getting all event types
func (h *Handler) GetEventTypes(c *gin.Context) {
	eventTypes, err := h.Service.GetEventTypes()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, eventTypes)
}

// GetEventType handles getting an event type by ID
func (h *Handler) GetEventType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	eventType, err := h.Service.GetEventTypeByID(id)
	if err != nil {
		respondWithError(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, eventType)
}

// UpdateEventType handles updating an event type
func (h *Handler) UpdateEventType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req models.UpdateEventTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	eventType, err := h.Service.UpdateEventType(id, &req)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, eventType)
}

// DeleteEventType handles deleting an event type
func (h *Handler) DeleteEventType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := h.Service.DeleteEventType(id); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event type deleted successfully"})
}

// CreateBadge handles creating a new badge
func (h *Handler) CreateBadge(c *gin.Context) {
	var req models.NewBadgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	badge, err := h.Service.CreateBadge(&req)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, badge)
}

// GetBadges handles getting all badges
func (h *Handler) GetBadges(c *gin.Context) {
	badges, err := h.Service.GetBadges()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, badges)
}

// GetActiveBadges handles getting all active badges
func (h *Handler) GetActiveBadges(c *gin.Context) {
	badges, err := h.Service.GetActiveBadges()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, badges)
}

// GetBadge handles getting a badge by ID
func (h *Handler) GetBadge(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	badge, err := h.Service.GetBadgeByID(id)
	if err != nil {
		respondWithError(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, badge)
}

// GetBadgeWithCriteria handles getting a badge with its criteria
func (h *Handler) GetBadgeWithCriteria(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	badge, err := h.Service.GetBadgeWithCriteria(id)
	if err != nil {
		respondWithError(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, badge)
}

// UpdateBadge handles updating a badge
func (h *Handler) UpdateBadge(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req models.UpdateBadgeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	badge, err := h.Service.UpdateBadge(id, &req)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, badge)
}

// DeleteBadge handles deleting a badge
func (h *Handler) DeleteBadge(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := h.Service.DeleteBadge(id); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Badge deleted successfully"})
}

// ProcessEvent handles processing an event
func (h *Handler) ProcessEvent(c *gin.Context) {
	var req models.NewEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.Service.ProcessEvent(&req); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event processed successfully"})
}

// GetUserBadges handles getting all badges awarded to a user
func (h *Handler) GetUserBadges(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		respondWithError(c, http.StatusBadRequest, "User ID is required")
		return
	}

	badges, err := h.Service.GetUserBadges(userID)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, badges)
}

// CreateConditionType handles creating a new condition type
func (h *Handler) CreateConditionType(c *gin.Context) {
	var req models.NewConditionTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	conditionType, err := h.Service.CreateConditionType(&req)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, conditionType)
}

// GetConditionTypes handles getting all condition types
func (h *Handler) GetConditionTypes(c *gin.Context) {
	conditionTypes, err := h.Service.GetConditionTypes()
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, conditionTypes)
}

// GetConditionType handles getting a condition type by ID
func (h *Handler) GetConditionType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	conditionType, err := h.Service.GetConditionTypeByID(id)
	if err != nil {
		respondWithError(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, conditionType)
}

// UpdateConditionType handles updating a condition type
func (h *Handler) UpdateConditionType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req models.UpdateConditionTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid request payload")
		return
	}

	conditionType, err := h.Service.UpdateConditionType(id, &req)
	if err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, conditionType)
}

// DeleteConditionType handles deleting a condition type
func (h *Handler) DeleteConditionType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		respondWithError(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	if err := h.Service.DeleteConditionType(id); err != nil {
		respondWithError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Condition type deleted successfully"})
}
