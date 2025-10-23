package auth

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IkBenJur/repetition-backend/config"
	"github.com/IkBenJur/repetition-backend/types"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func CreateJWT(secret []byte, userId int) (string, error) {
	tokenTimeout := time.Second * time.Duration(3600*24*7) // One week

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": strconv.Itoa(userId),
		"expiredAt": time.Now().Add(tokenTimeout).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(userController types.UserController) gin.HandlerFunc {
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

		userId, _ := strconv.Atoi(useIdString)

		if _, err := userController.GetUserById(userId); err != nil {
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