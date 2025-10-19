package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IkBenJur/repetition-backend/types"
	"github.com/gin-gonic/gin"
)

func TestUserControllersHandlers(t *testing.T) {
	userController := &mockUserController{}
	handler := NewHandler(userController)

	t.Run("Should fail if user payload is invalid", func(t *testing.T) {
		payload := types.RegisterUserPayload {
			Username: "NotJur",
			Password: "",
		}
		marshalled, _ := json.Marshal(payload)
		
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := gin.Default()

		router.POST("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status code %d, got %d", http.StatusBadRequest, rr.Code)
		}
	})
	t.Run("Create user when valid payload", func(t *testing.T) {
		payload := types.RegisterUserPayload {
			Username: "NotJur",
			Password: "a$$word",
		}
		marshalled, _ := json.Marshal(payload)
		
		req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(marshalled))
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		router := gin.Default()

		router.POST("/register", handler.handleRegister)
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("expected status code %d, got %d", http.StatusCreated, rr.Code)
		}
	})
}

type mockUserController struct {}

func (controller *mockUserController) GetUserByUsername(username string) (*types.User, error) {
	return nil, fmt.Errorf("Failed to find user")
}
func (controller *mockUserController) SaveUser(user types.User) error {
	return nil
}
func (controller *mockUserController) GetUserById(id int) (*types.User, error) {
	return nil, fmt.Errorf("Failed to find user")
}