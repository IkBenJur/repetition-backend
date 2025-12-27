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
	router.GET("/userWorkout", auth.WithJWTAuth(handler.userController), handler.handleGetAllUserWorkouts)
	router.GET("/userWorkout/active", auth.WithJWTAuth(handler.userController), handler.handleFindActiveUserWorkout)
	router.GET("/userWorkout/:id", auth.WithJWTAuth(handler.userController), handler.handleGetByUserWorkoutId)
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

	// Set the userId to that of the logged in user
	newUserWorkout.UserId = c.GetInt("userId")

	userWorkout := newUserWorkout.ToEntity()

	newWorkoutId, err := handler.controller.CreateNewUserWorkout(*userWorkout)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create"})
		return
	}

	err = handler.userController.UpdateActiveUserWorkoutForUserId(userWorkout.UserId, newWorkoutId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to set new active workout"})
		return
	}

	// Find the new workout from the database
	userWorkout, err = handler.controller.FindById(newWorkoutId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to find new workout"})
		return

	}

	c.JSON(http.StatusCreated, userWorkout)
}

func (handler *Handler) handleGetByUserWorkoutId(c *gin.Context) {
	userId := c.GetInt("userId")
	userWorkoutIdParam := c.Param("id")

	userWorkoutId, err := strconv.Atoi(userWorkoutIdParam)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	userWorkout, err := handler.controller.FindById(userWorkoutId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	if userId != userWorkout.UserId {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not your workout!"})
		return
	}

	c.JSON(http.StatusOK, userWorkout)
}

func (handler *Handler) handleGetAllUserWorkouts(c *gin.Context) {
	userId := c.GetInt("userId")

	userWorkouts, err := handler.controller.FindAllWorkoutsForUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, userWorkouts)
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

	// Find the userWorkout again
	userWorkout, err := handler.controller.FindById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	// Check if active workout for user should also be reset
	if user.ActiveUserWorkoutId == nil {
		c.JSON(http.StatusOK, gin.H{"userWorkout": userWorkout, "wasActiveWorkout": false})
		return
	}

	if *user.ActiveUserWorkoutId != id {
		c.JSON(http.StatusOK, gin.H{"userWorkout": userWorkout, "wasActiveWorkout": false})
		return
	}

	// Remove the active workout
	user.ActiveUserWorkoutId = nil
	err = handler.userController.UpdateUser(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"userWorkout": userWorkout, "wasActiveWorkout": true})
}
