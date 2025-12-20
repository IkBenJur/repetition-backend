package userWorkout

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/IkBenJur/repetition-backend/service/auth"
	"github.com/IkBenJur/repetition-backend/types"
	"github.com/IkBenJur/repetition-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	controller     Controller
	userController types.UserController
}

func NewHandler(controller Controller, userController types.UserController) *Handler {
	return &Handler{
		controller:     controller,
		userController: userController,
	}
}

func (handler *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/userWorkout", auth.WithJWTAuth(handler.userController), handler.handleCreateNewUserWorkout)
	router.GET("/userWorkout/active", auth.WithJWTAuth(handler.userController), handler.handleFindActiveUserWorkout)
	router.PUT("/userWorkout/:id/mark-complete", auth.WithJWTAuth(handler.userController), handler.handleMarkUserworkoutAsComplete)
}

func (handler *Handler) handleCreateNewUserWorkout(c *gin.Context) {
	var newUserWorkout types.NewUserWorkoutPayload

	if err := c.ShouldBindJSON(&newUserWorkout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	if err := utils.Validate.Struct(newUserWorkout); err != nil {
		errors := err.(validator.ValidationErrors)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid payload: %v", errors)})
		return
	}

	if newUserWorkout.UserId != c.GetInt("userId") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "haha dont even try it"})
		return
	}

	userWorkout := types.UserWorkoutPayloadIntoUserWorkout(newUserWorkout)

	newWorkoutId, err := handler.controller.CreateNewUserWorkout(userWorkout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create"})
		return
	}

	err = handler.userController.UpdateActiveUserWorkoutForUserId(userWorkout.UserId, newWorkoutId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set new active workout"})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (handler *Handler) handleFindActiveUserWorkout(c *gin.Context) {
	userId := c.GetInt("userId")

	userWorkout, err := handler.controller.FindActiveWorkoutForUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, userWorkout)
}

func (handler *Handler) handleMarkUserworkoutAsComplete(c *gin.Context) {
	idParam := c.Param("id")
	userId := c.GetInt("userId")

	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	workoutUserId, err := handler.controller.FindUserIdForUserworkoutId(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if workoutUserId != userId {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid"})
		return
	}

	err = handler.controller.MarkUserWorkoutAsComplete(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	user, err := handler.userController.GetUserById(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	// Check if active workout for user should also be reset
	if user.ActiveUserWorkoutId == nil {
		c.JSON(http.StatusOK, nil)
		return
	}

	if *user.ActiveUserWorkoutId != id {
		c.JSON(http.StatusOK, nil)
		return
	}

	// Remove the active workout
	user.ActiveUserWorkoutId = nil
	err = handler.userController.UpdateUser(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, nil)
}
