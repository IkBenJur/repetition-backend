package routes

import (
	"github.com/IkBenJur/repetition-backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/ping", controllers.GetPong)
}