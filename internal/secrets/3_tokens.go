package secrets

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/tomek7667/go-http-helpers/chii"
	"github.com/tomek7667/go-http-helpers/h"
	"github.com/tomek7667/go-http-helpers/utils"
	"github.com/tomek7667/secrets/internal/sqlc"
)

type CreateTokenDto struct {
	Token     string     `json:"token"`
	ExpiresAt *time.Time `json:"expires_at"`
}

type UpdateTokenDto struct {
	ExpiresAt *time.Time `json:"expires_at"`
}

func (s *Server) AddTokensRoutes() {
	auth := s.Router.With(chii.WithAuth(s.auther))
	auth.Route("/api/tokens", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			tokens, err := s.Db.Queries.ListTokens(r.Context())
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("failed to list tokens for user %s: %s", user.ID, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(GetTokensEvent, fmt.Sprintf("%s retrieved tokens", user.ID), r)
			}
			h.ResSuccess(w, tokens)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			dto, err := h.GetDto[CreateTokenDto](r)
			if err != nil {
				h.ResBadRequest(w, err)
				return
			}
			token, err := s.Db.Queries.CreateToken(r.Context(), sqlc.CreateTokenParams{
				ID:        utils.CreateUUID(),
				Token:     dto.Token,
				ExpiresAt: dto.ExpiresAt,
			})
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s tried to create token %s: %s", user.ID, dto.Token, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(IngestEvent, fmt.Sprintf("user %s created token %s", user.ID, token.ID), r)
			}
			h.ResSuccess(w, token)
		})

		r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			id := chi.URLParam(r, "id")
			dto, err := h.GetDto[UpdateTokenDto](r)
			if err != nil {
				h.ResBadRequest(w, err)
				return
			}
			token, err := s.Db.Queries.GetToken(r.Context(), id)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s tried to update token '%s' but an error happened: %s", user.ID, id, err.Error()), r)
				h.ResNotFound(w, "token")
				return
			}
			updatedToken, err := s.Db.Queries.UpdateToken(r.Context(), sqlc.UpdateTokenParams{
				ID:        id,
				ExpiresAt: dto.ExpiresAt,
			})
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to update token %s: %s", user.ID, id, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(UpdateTokenEvent, fmt.Sprintf("user %s from %s to %s", user.ID, token.ExpiresAt.Format(time.RFC3339), updatedToken.ExpiresAt.Format(time.RFC3339)), r)
			}
			h.ResSuccess(w, updatedToken)
		})

		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			id := chi.URLParam(r, "id")
			_, err := s.Db.Queries.GetToken(r.Context(), id)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to delete unexisting token %s: %s", user.ID, id, err.Error()), r)
				h.ResNotFound(w, "token")
				return
			}

			err = s.Db.Queries.DeleteToken(r.Context(), id)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to delete token %s: %s", user.ID, id, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(DeleteEvent, fmt.Sprintf("user %s deleted token %s", user.ID, id), r)
			}
			h.ResSuccess(w, nil)
		})
	})
}
