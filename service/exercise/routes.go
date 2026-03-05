package exercise

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/IkBenJur/repetition-backend/types"
	"github.com/IkBenJur/repetition-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ExerciseService interface {
	FindAll(ctx context.Context) ([]*types.Exercise, error)
	FindById(ctx context.Context, id int64) (*types.Exercise, error)
	Save(ctx context.Context, entity *types.Exercise) (int64, error)
}

type Handler struct {
	service ExerciseService
}

func NewHandler(service ExerciseService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/exercise", h.handleAllExercise)
	router.GET("/exercise/:id", h.handleGetExerciseById)
	router.POST("/exercise", h.handleNewExercise)
}

func (h *Handler) handleAllExercise(c *gin.Context) {
	exercises, err := h.service.FindAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, exercises)
}

func (h *Handler) handleNewExercise(c *gin.Context) {
	var payload types.ExercisePayload

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		errors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid payload: %v", errors)})
		return
	}

	entity, err := payload.ToEntity()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse entity"})
		return
	}

	newId, err := h.service.Save(c.Request.Context(), entity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
		return
	}

	entity.Id = &newId
	c.JSON(http.StatusCreated, entity)
}

func (h *Handler) handleGetExerciseById(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	exercise, err := h.service.FindById(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to find exercise"})
		return
	}

	c.JSON(http.StatusOK, exercise)
}

