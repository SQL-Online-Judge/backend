package restapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SQL-Online-Judge/backend/internal/model"
	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"go.uber.org/zap"
)

type loginResponse struct {
	Token string         `json:"token,omitempty"`
	Error *errorResponse `json:"error,omitempty"`
}

func (lr *loginResponse) toJSON() []byte {
	res, err := json.Marshal(lr)
	if err != nil {
		logger.Logger.Error("failed to marshal login response", zap.Error(err))
		return nil
	}
	return res
}

func login(w http.ResponseWriter, r *http.Request) {
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := loginResponse{Error: &errorResponse{Code: http.StatusBadRequest, Message: "invalid request"}}
		w.Write(resp.toJSON())
		return
	}

	if !user.IsValidLogin() {
		w.WriteHeader(http.StatusBadRequest)
		resp := loginResponse{Error: &errorResponse{Code: http.StatusBadRequest, Message: "invalid username or password"}}
		w.Write(resp.toJSON())
		return
	}

	if err := userService.CheckPassword(&user); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		resp := loginResponse{Error: &errorResponse{Code: http.StatusUnauthorized, Message: "invalid username or password"}}
		w.Write(resp.toJSON())
		return
	}

	userID := userService.GetUserIDByUsername(user.Username)
	if userID == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		resp := loginResponse{Error: &errorResponse{Code: http.StatusInternalServerError, Message: "internal server error"}}
		logger.Logger.Error("failed to get user id by username", zap.String("username", user.Username))
		w.Write(resp.toJSON())
		return
	}

	token, err := generateToken(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		resp := loginResponse{Error: &errorResponse{Code: http.StatusInternalServerError, Message: "internal server error"}}
		logger.Logger.Error("failed to generate token", zap.Error(err))
		w.Write(resp.toJSON())
		return
	}

	resp := loginResponse{Token: token, Error: nil}
	w.WriteHeader(http.StatusOK)
	w.Write(resp.toJSON())
}

func generateToken(userID int64) (string, error) {
	_, tokenString, err := tokenAuth.Encode(map[string]interface{}{
		"userID":          userID,
		jwt.ExpirationKey: time.Now().Add(time.Hour * 24),
	})
	if err != nil {
		return "", fmt.Errorf("failed to encode token: %w", err)
	}
	return tokenString, nil
}
