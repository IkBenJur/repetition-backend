package exercise

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/IkBenJur/repetition-backend/types"
	"github.com/IkBenJur/repetition-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	c types.ExerciseController
}

func NewHandler(c types.ExerciseController) *Handler {
	return &Handler{c: c}
}

func (h *Handler) RegisterRoutes(router *gin.Engine) {
	router.GET("/exercise", h.handleAllExercise)
	router.GET("/exercise/:id", h.handleGetExerciseById)
	router.POST("/exercise", h.handleNewExercise)
}

func (h *Handler) handleAllExercise(c *gin.Context) {
	exercises, err := h.c.GetAllExercise()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, exercises)
}

func (h *Handler) handleNewExercise(c *gin.Context) {
	var newExercise types.NewExercisePayload

	// Parse JSON
	if err := c.ShouldBindJSON(&newExercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Validate struct
	if err := utils.Validate.Struct(newExercise); err != nil {
		errors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid payload: %v", errors)})
		return
	}

	err := h.c.SaveExercise(types.Exercise{
		Name:        newExercise.Name,
		MuscleGroup: newExercise.MuscleGroup,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create exercise"})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (h *Handler) handleGetExerciseById(c *gin.Context) {
	idParam := c.Param("id")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	exercise, err := h.c.GetExerciseById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to find exercise"})
		return
	}

	c.JSON(http.StatusOK, exercise)
}
