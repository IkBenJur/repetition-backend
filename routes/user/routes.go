package user

import (
	"net/http"

	"github.com/IkBenJur/repetition-backend/controllers/auth"
	"github.com/IkBenJur/repetition-backend/types"
	"github.com/gin-gonic/gin"
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

}

func (handler *Handler) handleRegister(c *gin.Context) {
	var newUser types.RegisterUserPayload

	if err := c.ShouldBindJSON(&newUser); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
        return
    }
	
	// Check if user exists
	_, err := handler.controller.GetUserByUsername(newUser.Username)
	if err == nil {
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