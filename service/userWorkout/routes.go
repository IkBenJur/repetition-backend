package userWorkout

import (
	"fmt"
	"net/http"

	"github.com/IkBenJur/repetition-backend/types"
	"github.com/IkBenJur/repetition-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	controller Controller
	userController types.UserController
}

func NewHandler(controller Controller, userController types.UserController) *Handler {
	return &Handler{ 
		controller: controller,
		userController: userController,
	 }
}

func (handler *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/userWorkout", handler.handleSaveUserWorkout)
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

	if _, err := handler.userController.GetUserById(newUserWorkout.UserId); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("couldn't find user with ID: %v", newUserWorkout.UserId)})
		return
	}

	if err := handler.controller.SaveUserWorkout(types.UserWorkout{
		Name: newUserWorkout.Name,
		UserId: newUserWorkout.UserId,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create"})
	}

	c.JSON(http.StatusCreated, nil)
}