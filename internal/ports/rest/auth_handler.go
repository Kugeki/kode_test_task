package rest

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Kugeki/kode_test_task/internal/ports/rest/dto"
	"github.com/Kugeki/kode_test_task/internal/usecases"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"log/slog"
	"net/http"
	"time"
)

type AuthUsecase interface {
	VerifyLogin(ctx context.Context, username string, password string) error
}

type AuthHandler struct {
	uc AuthUsecase

	jwtExpireDuration time.Duration
	jwtSecretKey      string

	log *slog.Logger
}

func NewAuthHandler(log *slog.Logger, uc AuthUsecase, jwtSecretKey string, jwtExpireDuration time.Duration) *AuthHandler {
	h := &AuthHandler{
		uc:                uc,
		jwtExpireDuration: jwtExpireDuration,
		jwtSecretKey:      jwtSecretKey,
		log:               log.With(slog.String("context", "auth rest handler")),
	}
	return h
}

func (h *AuthHandler) SetupRoutes(r chi.Router) {
	r.Post("/users/login/", h.Login())
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := h.log.With(slog.String("request_id", middleware.GetReqID(ctx)))

		var req dto.LoginReq
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Warn("json decode error", slog.Any("error", err))
			respondError(w, http.StatusBadRequest, err, log)
			return
		}
		if req.Username == "" {
			log.Warn("username field is empty")
			respondError(w, http.StatusBadRequest, errors.New("username field is empty"), log)
			return
		}
		if req.Password == "" {
			log.Warn("password field is empty")
			respondError(w, http.StatusBadRequest, errors.New("password field is empty"), log)
			return
		}

		err = h.uc.VerifyLogin(r.Context(), req.Username, req.Password)
		if errors.Is(err, usecases.ErrWrongPassword) {
			log.Warn("verify login: unauthorized", slog.Any("error", err))
			respondUnauthorized(w, err, log)
			return
		}

		payload := jwt.MapClaims{
			"sub": req.Username,
			"exp": time.Now().Add(h.jwtExpireDuration).Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

		signedJWT, err := token.SignedString([]byte(h.jwtSecretKey))
		if err != nil {
			log.Error("jwt sign string error", slog.Any("error", err))
			respondError(w, http.StatusInternalServerError, err, log)
			return
		}

		resp := dto.LoginResp{AccessToken: signedJWT}
		respond(w, http.StatusOK, resp, log)
	}
}
