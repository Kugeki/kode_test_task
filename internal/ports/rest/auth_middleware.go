package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"strings"
)

var ContextJWTClaimsKey = "jwt claims"
var ContextUsernameKey = "jwt username"

type Middleware func(next http.Handler) http.Handler

func JWTAuthMiddleware(log *slog.Logger, jwtSecretKey string) Middleware {
	log = log.With(slog.String("middleware", "jwt auth middleware"))

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authField := r.Header.Get("Authorization")
			if !strings.HasPrefix(authField, "Bearer") {
				log.Warn("bearer prefix is not provided")
				respondUnauthorized(w, errors.New("need jwt token with Bearer prefix"), log)
				return
			}

			tokenString := strings.TrimPrefix(authField, "Bearer")
			tokenString = strings.Trim(tokenString, " ")

			token, err := parseToken(tokenString, jwtSecretKey)
			switch {
			case errors.Is(err, jwt.ErrTokenMalformed) || errors.Is(err, jwt.ErrTokenSignatureInvalid):
				log.Warn("parse token error", slog.Any("error", err))
				respondError(w, http.StatusBadRequest, err, log)
				return
			case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
				log.Warn("parse token error", slog.Any("error", err))
				respondUnauthorized(w, err, log)
				return
			case err != nil:
				log.Error("parse token error", slog.Any("error", err))
				respondError(w, http.StatusInternalServerError, err, log)
				return
			}

			ctx := context.WithValue(r.Context(), ContextJWTClaimsKey, token.Claims)

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Error("can't cast jwt token to map claims")
				respondError(w, http.StatusInternalServerError, errors.New("can't cast jwt token to map claims"), log)
				return
			}

			username, ok := claims["sub"].(string)
			if !ok {
				log.Error("can't get nickname from claims")
				respondError(w, http.StatusInternalServerError, errors.New("can't get nickname from claims"), log)
				return
			}
			ctx = context.WithValue(ctx, ContextUsernameKey, username)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func parseToken(tokenString, jwtSecretKey string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(jwtSecretKey), nil
	})
}
