package userWorkoutExerciseSet

import (
	"fmt"
	"net/http"

	userWorkoutExercise "github.com/IkBenJur/repetition-backend/service/UserWorkoutExercise"
	"github.com/IkBenJur/repetition-backend/service/authMiddleware"
	"github.com/IkBenJur/repetition-backend/service/user"
	types "github.com/IkBenJur/repetition-backend/types/userWorkout"
	"github.com/IkBenJur/repetition-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	controller                    UserWorkoutExerciseSetController
	userController                user.UserController
	userWorkoutExerciseController userWorkoutExercise.UserWorkoutExerciseController
}

func NewHandler(controller UserWorkoutExerciseSetController, userController user.UserController, userWorkoutExerciseController userWorkoutExercise.UserWorkoutExerciseController) *Handler {
	return &Handler{
		controller:                    controller,
		userController:                userController,
		userWorkoutExerciseController: userWorkoutExerciseController,
	}
}

func (handler *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/userWorkoutExerciseSet", authMiddleware.WithJWTAuth(handler.userController), handler.handleCreateOrUpdateUserWorkoutExerciseSet)
	router.PUT("/userWorkoutExerciseSet/:id", authMiddleware.WithJWTAuth(handler.userController), handler.handleCreateOrUpdateUserWorkoutExerciseSet)
}

func (handler *Handler) handleCreateOrUpdateUserWorkoutExerciseSet(c *gin.Context) {
	userId := c.GetInt64("userId")
	if userId == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var payload types.UserWorkoutExerciseSetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	if err := utils.Validate.Struct(payload); err != nil {
		validationErrors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Validation failed: %v", validationErrors)})
		return
	}

	userWorkoutExerciseSet, err := payload.ToEntity()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse entity"})
		return
	}

	// Check if should update or create
	if payload.IsUpdate() {
		// UPDATE existing set

		// Additional authorization check: verify user owns the set being updated
		if userId != *payload.UserId {
			c.JSON(http.StatusForbidden, gin.H{"error": "Not your set"})
			return
		}

		if _, err := handler.controller.Save(c.Request.Context(), userWorkoutExerciseSet); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update set"})
			return
		}

		c.JSON(http.StatusOK, userWorkoutExerciseSet)
	} else {
		// CREATE new set
		workoutUserId, err := handler.userWorkoutExerciseController.FindUserIdForUserWorkoutExerciseId(c.Request.Context(), payload.UserWorkoutExerciseId)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Workout exercise not found"})
			return
		}

		// Validate user own workout
		if workoutUserId != userId {
			c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to modify this workout"})
			return
		}

		// Determine the new setNumber
		setNumber, err := handler.controller.DetermineSetNumberForNewUserWorkoutExerciseSet(c.Request.Context(), userWorkoutExerciseSet.UserWorkoutExerciseId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to determine set number"})
			return
		}

		userWorkoutExerciseSet.SetNumber = &setNumber

		userWorkoutSetId, err := handler.controller.Save(c.Request.Context(), userWorkoutExerciseSet)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create set"})
			return
		}

		userWorkoutExerciseSet.Id = &userWorkoutSetId
		c.JSON(http.StatusCreated, userWorkoutExerciseSet)
	}
}
