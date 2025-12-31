package workouttemplate

import (
	"fmt"
	"net/http"

	"github.com/IkBenJur/repetition-backend/service/auth"
	"github.com/IkBenJur/repetition-backend/types"
	"github.com/IkBenJur/repetition-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	Controller     Controller
	userController types.UserController
}

func NewHandler(controller Controller, userController types.UserController) *Handler {
	return &Handler{
		Controller:     controller,
		userController: userController,
	}
}

func (handler *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/workoutTemplate", auth.WithJWTAuth(handler.userController))
}

// TODO For now only create but will handle update later
func (handler *Handler) handleCreateOrUpdateNewWorkout(c *gin.Context) {
	var payload types.TemplateWorkoutPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Validation failed: %v", validationErrors)})
		return
	}

	userId := c.GetInt("userId")
	templateWorkout := payload.ToEntity()

	templateWorkout.UserId = userId

	if _, err := handler.Controller.CreateNewTemplateWorkout(templateWorkout); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create template"})
		return
	}

	c.JSON(http.StatusCreated, templateWorkout)
}
