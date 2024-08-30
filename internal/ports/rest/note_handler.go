package rest

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/Kugeki/kode_test_task/internal/domain"
	"github.com/Kugeki/kode_test_task/internal/ports/rest/dto"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log/slog"
	"net/http"
)

type NoteUsecase interface {
	ValidateNoteSpelling(ctx context.Context, n *domain.Note) (domain.SpellResults, error)
	CreateNote(ctx context.Context, authorName string, n *domain.Note) error
	GetNotesForUser(ctx context.Context, username string) ([]*domain.Note, error)
}

type NoteHandler struct {
	uc           NoteUsecase
	jwtSecretKey string

	log *slog.Logger
}

func NewNoteHandler(log *slog.Logger, uc NoteUsecase, jwtSecretKey string) *NoteHandler {
	h := &NoteHandler{
		uc:           uc,
		jwtSecretKey: jwtSecretKey,
		log:          log.With(slog.String("context", "note rest handler")),
	}
	return h
}

func (h *NoteHandler) SetupRoutes(r chi.Router) {
	r.Group(func(r chi.Router) {
		r.Use(JWTAuthMiddleware(h.log, h.jwtSecretKey))

		r.Post("/notes/create/", h.CreateNote())
		r.Get("/notes/", h.GetNotes())
	})
}

// CreateNote godoc
//
//	@Summary		Create a note
//	@Description	create a note for user
//	@Tags			notes
//	@Accept			json
//	@Produce		json
//	@Param			note			body		dto.CreateNoteReq	true	"Create note"
//	@Param			Authorization	header		string				true	"Auth JWT Token"	default(Bearer <Add access token here>)
//	@Success		201				{object}	dto.CreateNoteResp
//	@Failure		400				{object}	dto.NoteSpellErrorResp
//	@Failure		401				{object}	dto.HTTPError
//	@Failure		500				{object}	dto.HTTPError
//
//	@Header			401				{string}	WWW-Authenticate	"Auth realm"
//
//	@Router			/notes/create/ [post]
//	@Security		BearerAuth
func (h *NoteHandler) CreateNote() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := h.log.With(slog.String("request_id", middleware.GetReqID(ctx)))

		req := dto.CreateNoteReq{}
		err := json.NewDecoder(r.Body).Decode(&req)
		if errors.Is(err, io.EOF) {
			log.Warn("no json body", slog.Any("error", err))
			respondError(w, http.StatusBadRequest, errors.New("need json body with content"), log)
			return
		}
		if err != nil {
			log.Error("json decode error", slog.Any("error", err))
			respondError(w, http.StatusInternalServerError, err, log)
			return
		}
		if req.Content == "" {
			log.Warn("note: content field is empty")
			respondError(w, http.StatusBadRequest, errors.New("note: content field is empty"), log)
			return
		}

		note := req.ToDomain()

		spellResults, err := h.uc.ValidateNoteSpelling(ctx, note)
		if err != nil {
			log.Error("validate note spelling error", slog.Any("error", err))
			respondError(w, http.StatusInternalServerError, err, log)
			return
		}
		if len(spellResults) > 0 {
			resp := &dto.NoteSpellErrorResp{}
			resp.FromDomain(note, spellResults)
			respond(w, http.StatusBadRequest, resp, log)
			return
		}

		username, ok := ctx.Value(ContextUsernameKey).(string)
		if !ok {
			log.Error("can't get nickname from context")
			respondError(w, http.StatusInternalServerError, errors.New("can't get nickname from context"), log)
			return
		}

		err = h.uc.CreateNote(ctx, username, note)
		if err != nil {
			log.Error("create note error", slog.Any("error", err))
			respondError(w, http.StatusInternalServerError, err, log)
			return
		}

		resp := &dto.CreateNoteResp{}
		resp.FromDomain(note)
		respond(w, http.StatusCreated, resp, log)
	}
}

// GetNotes godoc
//
//	@Summary		Get notes
//	@Description	get notes for user
//	@Tags			notes
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Auth JWT Token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	dto.GetNotesResp
//	@Failure		400				{object}	dto.HTTPError
//	@Failure		401				{object}	dto.HTTPError
//	@Failure		500				{object}	dto.HTTPError
//
//	@Header			401				{string}	WWW-Authenticate	"Auth realm"
//
//	@Router			/notes/ [get]
//	@Security		BearerAuth
func (h *NoteHandler) GetNotes() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log := h.log.With(slog.String("request_id", middleware.GetReqID(ctx)))

		username, ok := ctx.Value(ContextUsernameKey).(string)
		if !ok {
			log.Error("can't get nickname from context")
			respondError(w, http.StatusInternalServerError, errors.New("can't get nickname from context"), log)
			return
		}

		notes, err := h.uc.GetNotesForUser(ctx, username)
		if err != nil {
			log.Error("get notes for user error", slog.Any("error", err))
			respondError(w, http.StatusInternalServerError, err, log)
			return
		}

		resp := &dto.GetNotesResp{}
		resp.FromDomain(notes)
		respond(w, http.StatusOK, resp, log)
	}
}
