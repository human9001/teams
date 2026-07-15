package v1

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/human9001/teams/internal/infrastructure/auth"
	apiErrors "github.com/human9001/teams/internal/interfaces/http/api/v1/errors"
	"github.com/human9001/teams/pkg/helpers"
)

type ctxKey string

const ClaimsKey ctxKey = "jwt_claims"

func AuthMiddleware(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
				helpers.WriteError(w, http.StatusUnauthorized, apiErrors.ErrBearerToken)
				return
			}

			tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.ParseWithClaims(tokenStr, &auth.Claims{}, func(token *jwt.Token) (any, error) {
				return []byte(secret), nil
			})
			if err != nil || !token.Valid {
				helpers.WriteError(w, http.StatusUnauthorized, apiErrors.ErrInvalidToken)
				return
			}

			claims, ok := token.Claims.(*auth.Claims)
			if !ok {
				helpers.WriteError(w, http.StatusUnauthorized, apiErrors.ErrInvalidTokenClaims)
				return
			}

			ctx := context.WithValue(r.Context(), ClaimsKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func ClaimsFromContext(ctx context.Context) (*auth.Claims, bool) {
	claims, ok := ctx.Value(ClaimsKey).(*auth.Claims)
	return claims, ok
}
