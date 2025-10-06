package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetPong(c *gin.Context) {
	c.JSON(http.StatusOK, "pong")
}