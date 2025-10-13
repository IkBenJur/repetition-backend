package user

import (
	"fmt"
	"net/http"

	"github.com/IkBenJur/repetition-backend/controllers/auth"
	"github.com/IkBenJur/repetition-backend/types"
	"github.com/IkBenJur/repetition-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	controller types.UserController
}

func NewHandler(controller types.UserController) *Handler {
	return &Handler{ controller: controller }
}

func (handler *Handler) RegisterRoutes(router *gin.Engine) {
	router.POST("/login", handler.handleLogin)
	router.POST("/register", handler.handleRegister)
}

func (handler *Handler) handleLogin(c *gin.Context) {
	var loginUser types.LoginUserPayload

	// Parse JSON
	if err := c.ShouldBindJSON(&loginUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

	// Validate struct
	if err := utils.Validate.Struct(loginUser); err != nil {
		errors := err.(validator.ValidationErrors)
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid payload: %v", errors)})
		return
	}

	user, err := handler.controller.GetUserByUsername(loginUser.Username)
	if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login"})
		return
	}

	if !auth.ComparePassword(user.Password, loginUser.Password) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": ""})
}

func (handler *Handler) handleRegister(c *gin.Context) {
	var newUser types.RegisterUserPayload

	// Parse JSON
	if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }

	// Validate struct
	if err := utils.Validate.Struct(newUser); err != nil {
		errors := err.(validator.ValidationErrors)
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid payload: %v", errors)})
		return
	}
	
	// Check if user exists
	_, err := handler.controller.GetUserByUsername(newUser.Username)
	if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User with username already exists"})
		return
	}
	
	// Hash
	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	
	// Create new user
	err = handler.controller.SaveUser(types.User{
		Username: newUser.Username,
		Password: hashedPassword,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, nil)
}