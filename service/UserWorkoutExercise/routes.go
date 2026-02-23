package userWorkoutExercise

import (
	"fmt"
	"net/http"

	"github.com/IkBenJur/repetition-backend/service/authMiddleware"
	"github.com/IkBenJur/repetition-backend/service/user"
	"github.com/IkBenJur/repetition-backend/service/userWorkout"
	types "github.com/IkBenJur/repetition-backend/types/userWorkout"
	"github.com/IkBenJur/repetition-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	controller            UserWorkoutExerciseController
	userController        user.UserController
	userWorkoutController userWorkout.Controller
}

func NewHandler(controller UserWorkoutExerciseController, userController user.UserController, userWorkoutController userWorkout.Controller) *Handler {
	return &Handler{
		controller:            controller,
		userController:        userController,
		userWorkoutController: userWorkoutController,
	}
}

func (handler *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/userWorkoutExercise", authMiddleware.WithJWTAuth(handler.userController), handler.handleCreateNewUserWorkoutExercise)
}

func (handler *Handler) handleCreateNewUserWorkoutExercise(c *gin.Context) {
	var newUserWorkoutExercise types.UserWorkoutExercisePayload

	if err := c.ShouldBindJSON(&newUserWorkoutExercise); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := utils.Validate.Struct(newUserWorkoutExercise); err != nil {
		errors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid payload: %v", errors)})
		return
	}

	// Check if user is allowed to add the exercise
	workoutUserId, err := handler.userWorkoutController.FindUserIdForUserworkoutId(newUserWorkoutExercise.UserWorkoutId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	if workoutUserId != c.GetInt("userId") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userWorkoutExercise, err := newUserWorkoutExercise.ToEntity()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse entity"})
		return
	}

	exerciseNumber, err := handler.controller.DetermineExerciseNumberForNewUserWorkoutExercise(c.Request.Context(), userWorkoutExercise.UserWorkoutId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to determine exercise number"})
		return
	}

	userWorkoutExercise.ExerciseNumber = &exerciseNumber

	userWorkoutExerciseId, err := handler.controller.Save(c.Request.Context(), userWorkoutExercise)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create"})
		return
	}

	userWorkoutExercise.Id = &userWorkoutExerciseId

	c.JSON(http.StatusCreated, userWorkoutExercise)
}
