package exercise

import (
	"net/http"

	"github.com/IkBenJur/repetition-backend/types"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	c types.ExerciseController
}

func NewHandler(c types.ExerciseController) *Handler {
	return &Handler{ c: c}
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

	c.JSON(http.StatusOK, gin.H{"exercises": exercises})
}

func (h *Handler) handleNewExercise(c *gin.Context) {

}

func (h *Handler) handleGetExerciseById(c *gin.Context) {

}