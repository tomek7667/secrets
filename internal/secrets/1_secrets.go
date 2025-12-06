package secrets

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/tomek7667/go-http-helpers/chii"
	"github.com/tomek7667/go-http-helpers/h"
	"github.com/tomek7667/go-http-helpers/utils"
	"github.com/tomek7667/secrets/internal/sqlc"
)

type CreateSecretDto struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type UpdateSecretDto struct {
	Value string `json:"value"`
}

func getSupportedTokenTypes() []string {
	return []string{"Api"}
}

func (s *Server) AddSecretsRoutes() {
	auth := s.Router.With(chii.WithAuth(s.auther))
	auth.Route("/api/secrets", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			secrets, err := s.Db.Queries.ListSecrets(r.Context())
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("failed to list secrets for user %s: %s", user.ID, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(GetSecretsEvent, fmt.Sprintf("%s retrieved secrets", user.ID), r)
			}
			h.ResSuccess(w, secrets)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			dto, err := h.GetDto[CreateSecretDto](r)
			if err != nil {
				h.ResBadRequest(w, err)
				return
			}
			secret, err := s.Db.Queries.CreateSecret(r.Context(), sqlc.CreateSecretParams{
				ID:    utils.CreateUUID(),
				Key:   dto.Key,
				Value: utils.B64Encode(dto.Value),
			})
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s tried to create secret %s: %s", user.ID, dto.Key, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(IngestEvent, fmt.Sprintf("user %s created secret %s", user.ID, dto.Key), r)
			}
			h.ResSuccess(w, secret)
		})

		r.Put("/", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			key := r.URL.Query().Get("key")
			dto, err := h.GetDto[UpdateSecretDto](r)
			if err != nil {
				h.ResBadRequest(w, err)
				return
			}
			secret, err := s.Db.Queries.GetSecret(r.Context(), key)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s tried to update secret '%s' but an error happened: %s", user.ID, key, err.Error()), r)
				h.ResNotFound(w, "secret")
				return
			}
			updatedSecret, err := s.Db.Queries.UpdateSecret(r.Context(), sqlc.UpdateSecretParams{
				Key:   key,
				Value: dto.Value,
			})
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to update secret %s: %s", user.ID, key, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(UpdateSecretEvent, fmt.Sprintf("user %s from %s to %s", user.ID, secret.Value, updatedSecret.Value), r)
			}
			h.ResSuccess(w, updatedSecret)
		})

		r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			key := r.URL.Query().Get("key")

			_, err := s.Db.Queries.GetSecret(r.Context(), key)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to delete unexisting secret %s: %s", user.ID, key, err.Error()), r)
				h.ResNotFound(w, "secret")
				return
			}

			err = s.Db.Queries.DeleteSecret(r.Context(), key)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to delete secret %s: %s", user.ID, key, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(DeleteEvent, fmt.Sprintf("user %s deleted secret %s", user.ID, key), r)
			}
			h.ResSuccess(w, nil)
		})
	})

	s.Router.Get("/api/secrets/get", func(w http.ResponseWriter, r *http.Request) {
		var err error
		var tkn sqlc.Token
		authValue := strings.TrimSpace(r.Header.Get("Authorization"))
		key := r.URL.Query().Get("key")
		if strings.HasPrefix(authValue, "Api ") {
			token, _ := strings.CutPrefix(authValue, "Api ")
			tkn, err = s.Db.Queries.GetToken(r.Context(), token)
			if err != nil {
				s.Log(UnauthorizedEvent, fmt.Sprintf("invalid token '%s' provided to get secret '%s': %s", authValue, key, err.Error()), r)
				h.ResUnauthorized(w)
				return
			}
		} else {
			s.Log(UnauthorizedEvent, fmt.Sprintf("invalid token type '%s' provided to get secret '%s'. Supported token types: %s", authValue, key, strings.Join(getSupportedTokenTypes(), ", ")), r)
			h.ResUnauthorized(w)
			return
		}
		permissions, err := s.Db.Queries.ListPermissionsByTokenId(r.Context(), tkn.ID)
		if err != nil {
			s.Log(ErrorEvent, fmt.Sprintf("failed to list permissions for token %s: %s", tkn.ID, err.Error()), r)
			h.ResErr(w, err)
			return
		}
		secret, err := s.Db.Queries.GetSecret(r.Context(), key)
		if err != nil {
			s.Log(ErrorEvent, fmt.Sprintf("(before matching) couldn't retrieve secret %s for token %s: %s", tkn.ID, key, err.Error()), r)
			h.ResErr(w, err)
			return
		}
		matched := false
		for _, permission := range permissions {
			if PatternMatches(secret.Key, permission.SecretKeyPattern) {
				matched = true
				break
			}
		}
		if !matched {
			s.Log(UnauthorizedEvent, fmt.Sprintf("token %s can't access %s", tkn.ID, key), r)
			h.ResUnauthorized(w)
			return
		}
		s.Log(GetSecretEvent, fmt.Sprintf("token %s", tkn.ID), r)
		h.ResSuccess(w, secret)
	})
}
