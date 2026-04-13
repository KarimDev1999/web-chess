package middleware

import (
	"context"
	"net/http"
	"strings"

	"chess-backend/internal/constants/appconst"
	"chess-backend/pkg/jwt"
)

type contextKey string

const UserIDKey contextKey = "userID"

func AuthMiddleware(tokenSvc *jwt.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get(appconst.HeaderAuthorization)
			if authHeader == "" {
				http.Error(w, appconst.MsgMissingAuthHeader, http.StatusUnauthorized)
				return
			}
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != appconst.AuthSchemeBearer {
				http.Error(w, appconst.MsgInvalidAuthHeader, http.StatusUnauthorized)
				return
			}
			userID, err := tokenSvc.Validate(parts[1])
			if err != nil {
				http.Error(w, appconst.MsgInvalidToken, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (string, bool) {
	val := ctx.Value(UserIDKey)
	if val == nil {
		return "", false
	}
	userID, ok := val.(string)
	return userID, ok
}

func WSAuthMiddleware(tokenSvc *jwt.TokenService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := r.URL.Query().Get(appconst.QueryParamToken)
			if token == "" {
				http.Error(w, appconst.MsgMissingToken, http.StatusUnauthorized)
				return
			}
			userID, err := tokenSvc.Validate(token)
			if err != nil {
				http.Error(w, appconst.MsgInvalidToken, http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
