package userWorkout

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
	router.POST("/userWorkout", auth.WithJWTAuth(handler.userController), handler.handleSaveUserWorkout)
	router.GET("/active-userWorkout", auth.WithJWTAuth(handler.userController), handler.handleFindActiveUserWorkout)
}

func (handler *Handler) handleSaveUserWorkout(c *gin.Context) {
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
	
	if err := handler.controller.SaveUserWorkout(userWorkout); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create"})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

func (handler *Handler) handleFindActiveUserWorkout(c *gin.Context) {
	userId := c.GetInt("userId")

	fmt.Printf("%v", userId)

	userWorkout, err := handler.controller.FindActiveWorkoutForUserId(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	c.JSON(http.StatusFound, userWorkout)
}

