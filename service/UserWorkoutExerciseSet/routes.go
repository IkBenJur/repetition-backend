package userWorkoutExerciseSet

import (
	"fmt"
	"net/http"

	userWorkoutExercise "github.com/IkBenJur/repetition-backend/service/UserWorkoutExercise"
	"github.com/IkBenJur/repetition-backend/service/auth"
	"github.com/IkBenJur/repetition-backend/types"
	"github.com/IkBenJur/repetition-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	controller                    Controller
	userController                types.UserController
	userWorkoutExerciseController userWorkoutExercise.Controller
}

func NewHandler(controller Controller, userController types.UserController, userWorkoutExerciseController userWorkoutExercise.Controller) *Handler {
	return &Handler{
		controller:                    controller,
		userController:                userController,
		userWorkoutExerciseController: userWorkoutExerciseController,
	}
}

func (handler *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/userWorkoutExerciseSet", auth.WithJWTAuth(handler.userController), handler.handleCreateNewUserWorkoutExerciseSet)
}

func (handler *Handler) handleCreateNewUserWorkoutExerciseSet(c *gin.Context) {
	var newUserWorkoutExerciseSet types.UserWorkoutExerciseSetPayload

	if err := c.ShouldBindJSON(&newUserWorkoutExerciseSet); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
		return
	}

	if err := utils.Validate.Struct(newUserWorkoutExerciseSet); err != nil {
		errors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid payload: %v", errors)})
		return
	}

	// Check if user is allowed to add the set to exercise
	workoutUserId, err := handler.userWorkoutExerciseController.FindUserIdForUserWorkoutExerciseId(newUserWorkoutExerciseSet.UserWorkoutExerciseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if workoutUserId != c.GetInt("userId") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userWorkoutExerciseSet := types.UserWorkoutExerciseSetIntoUserWorkoutExerciseSet(newUserWorkoutExerciseSet)
	userWorkoutSetId, err := handler.controller.CreateNewUserWorkoutExerciseSet(*userWorkoutExerciseSet)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create"})
		return
	}

	// Set the new ID
	userWorkoutExerciseSet.ID = userWorkoutSetId

	c.JSON(http.StatusCreated, userWorkoutExerciseSet)
}
