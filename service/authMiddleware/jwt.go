package authMiddleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/IkBenJur/repetition-backend/config"
	"github.com/IkBenJur/repetition-backend/service/user"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func WithJWTAuth(userController user.UserController) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := getJwtToken(c.Request)
		if tokenString == "" {
			log.Printf("No token present")
			permissionDenied(c)
			return
		}
		token, err := validateToken(tokenString)
		if err != nil {
			log.Printf("Failed to validate token: %v", err)
			permissionDenied(c)
			return
		}

		if !token.Valid {
			log.Printf("Token invalid: %v", token)
			permissionDenied(c)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		useIdString := claims["userId"].(string)

		userId, _ := strconv.ParseInt(useIdString, 10, 64)

		if _, err := userController.FindById(context.Background(), int64(userId)); err != nil {
			log.Printf("User does not exists: %v", userId)
			permissionDenied(c)
			return
		}

		// TODO Check token expires
		c.Set("userId", userId)

	}
}

func getJwtToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	return token
}

func validateToken(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})
}

func permissionDenied(c *gin.Context) {
	c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "permision denied"})
}
