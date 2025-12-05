package secrets

import (
	"fmt"
	"net/http"

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

func (s *Server) AddSecretsRoutes() {
	auth := s.Router.With(chii.WithAuth(s.auther))
	auth.Route("/api/secrets", func(r chi.Router) {
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
				s.Log(ErrorEvent, fmt.Sprintf("user %s tried to create secret %s: %s", user.ID, dto.Key, err.Error()))
				h.ResErr(w, err)
				return
			} else {
				s.Log(IngestEvent, fmt.Sprintf("user %s created secret %s", user.ID, dto.Key))
			}
			h.ResSuccess(w, secret)
		})

		r.Get("/get", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			key := r.URL.Query().Get("key")
			secret, err := s.Db.Queries.GetSecret(r.Context(), key)
			if err != nil {
				s.Log(GetSecretEvent, fmt.Sprintf("user %s tried to access secret '%s' but an error happened: %s", user.ID, key, err.Error()))
				h.ResNotFound(w, "secret")
				return
			} else {
				s.Log(GetSecretEvent, fmt.Sprintf("user %s", user.ID))
			}
			h.ResSuccess(w, secret)
		})
	})
}
