package main

import (
	"github.com/IkBenJur/repetition-backend/config"
	"github.com/IkBenJur/repetition-backend/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDatabase()

	router := gin.Default()
	routes.SetupRoutes(router)

	router.Run(":8080")
}
