package restapi

import (
	"context"
	"net/http"
	"strconv"

	"github.com/SQL-Online-Judge/backend/internal/pkg/logger"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"go.uber.org/zap"
)

func contentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func getRequestID(r *http.Request) string {
	requestID, ok := r.Context().Value(middleware.RequestIDKey).(string)
	if !ok {
		logger.Logger.Error("failed to get request ID from context")
		return ""
	}
	return requestID
}

type userIDKeyType struct{}
type roleKeyType struct{}

var userIDKey = userIDKeyType{}
var roleKey = roleKeyType{}

func getRole(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := getRequestID(r)
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil {
			logger.Logger.Error("failed to get claims from context", zap.String("requestID", requestID), zap.Error(err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		iUserID, ok := claims["userID"]
		if !ok {
			logger.Logger.Error("failed to get userID from claims", zap.String("requestID", requestID))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		sUserID, ok := iUserID.(string)
		if !ok {
			logger.Logger.Error("failed to convert userID to string", zap.String("requestID", requestID))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		userID, err := strconv.ParseInt(sUserID, 10, 64)
		if err != nil {
			logger.Logger.Error("failed to parse userID", zap.String("requestID", requestID), zap.Error(err))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		role := userService.GetRoleByUserID(userID)
		if role == "" {
			logger.Logger.Error("failed to get role by userID", zap.String("requestID", requestID), zap.Int64("userID", userID))
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		ctx = context.WithValue(ctx, roleKey, role)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func checkRole(acceptedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := getRequestID(r)
			role, ok := r.Context().Value(roleKey).(string)
			if !ok {
				logger.Logger.Error("failed to get role from context", zap.String("requestID", requestID))
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			for _, acceptedRole := range acceptedRoles {
				if role == acceptedRole {
					next.ServeHTTP(w, r)
					return
				}
			}

			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		})
	}
}
