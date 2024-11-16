package middlewares

import (
	"context"
	"errors"
	"fmt"
	"github/devAshu12/learning_platform_GO_backend/internal/auth"
	"github/devAshu12/learning_platform_GO_backend/internal/utils"
	"github/devAshu12/learning_platform_GO_backend/pkg/db"
	"github/devAshu12/learning_platform_GO_backend/pkg/models"
	"github/devAshu12/learning_platform_GO_backend/pkg/types"
	"net/http"

	"gorm.io/gorm"
)

type RequestUser struct {
	ID   string
	Role string
}

type UserContextKey string

const UserKey UserContextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		access_cookie, err := r.Cookie("access_token")
		if err != nil {
			appErr := types.NewAppError(http.StatusUnauthorized, "Missing or Invalid Auth Cookie", err)
			utils.RespondWithError(w, appErr)
			return
		}

		access_token := access_cookie.Value
		if access_token == "" {
			appErr := types.NewAppError(http.StatusUnauthorized, "Invalid Token", err)
			utils.RespondWithError(w, appErr)
			return
		}

		claims, err := auth.ValidateAccessToken(access_token, true)
		if err != nil {
			appErr := types.NewAppError(http.StatusUnauthorized, "Invalid Token", err)
			utils.RespondWithError(w, appErr)
			return
		}

		var user models.User
		if err := db.DB.First(&user, "id = ?", claims.Subject).Error; errors.Is(err, gorm.ErrRecordNotFound) {
			appErr := types.NewAppError(http.StatusUnauthorized, "Invalid Token", err)
			utils.RespondWithError(w, appErr)
			return
		} else if err != nil {
			appErr := types.NewAppError(http.StatusInternalServerError, "Internal error", err)
			utils.RespondWithError(w, appErr)
			return
		}

		requestUser := &RequestUser{
			ID:   user.ID,
			Role: string(user.Role),
		}

		ctx := context.WithValue(r.Context(), UserKey, requestUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(r *http.Request) (*RequestUser, error) {
	user, ok := r.Context().Value(UserKey).(*RequestUser)
	fmt.Print(user)
	if !ok || user == nil {
		return nil, errors.New("user not found in context")
	}
	return user, nil
}
