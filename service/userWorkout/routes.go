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
}

func (handler *Handler) handleSaveUserWorkout(c *gin.Context) {
	var newUserWorkout types.NewUserWorkoutPayload

	if err := c.ShouldBindJSON(&newUserWorkout); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}
	fmt.Print(newUserWorkout.UserWorkoutExercises)

	for _, test := range newUserWorkout.UserWorkoutExercises {
		fmt.Println(test.ExerciseId)
		fmt.Println(test.UserWorkoutExerciseSets)
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

	userWorkout := UserWorkoutPayloadIntoUserWorkout(newUserWorkout)
	fmt.Print(userWorkout)
	if err := handler.controller.SaveUserWorkout(userWorkout); err != nil {
		fmt.Printf("Error occured %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create"})
		return
	}

	c.JSON(http.StatusCreated, nil)
}

