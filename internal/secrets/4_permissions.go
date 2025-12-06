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

type CreatePermissionDto struct {
	SecretKeyPattern string `json:"secret_key_pattern"`
	TokenID          string `json:"token_id"`
}

type UpdatePermissionDto struct {
	SecretKeyPattern string `json:"secret_key_pattern"`
}

func (s *Server) AddPermissionsRoutes() {
	auth := s.Router.With(chii.WithAuth(s.auther))
	auth.Route("/api/permissions", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			permissions, err := s.Db.Queries.ListPermissions(r.Context())
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("failed to list permissions for user %s: %s", user.ID, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(GetPermissionsEvent, fmt.Sprintf("%s retrieved permissions", user.ID), r)
			}
			h.ResSuccess(w, permissions)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			dto, err := h.GetDto[CreatePermissionDto](r)
			if err != nil {
				h.ResBadRequest(w, err)
				return
			}
			_, err = s.Db.Queries.GetToken(r.Context(), dto.TokenID)
			if err != nil {
				h.ResNotFound(w, "specified token")
				return
			}
			permission, err := s.Db.Queries.CreatePermission(r.Context(), sqlc.CreatePermissionParams{
				ID:               utils.CreateUUID(),
				TokenID:          dto.TokenID,
				SecretKeyPattern: dto.SecretKeyPattern,
			})
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s tried to create permission %s for token %s: %s", user.ID, dto.SecretKeyPattern, dto.TokenID, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(IngestEvent, fmt.Sprintf("user %s created permission %s", user.ID, permission.ID), r)
			}
			h.ResSuccess(w, permission)
		})

		r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			id := chi.URLParam(r, "id")
			dto, err := h.GetDto[UpdatePermissionDto](r)
			if err != nil {
				h.ResBadRequest(w, err)
				return
			}
			permission, err := s.Db.Queries.GetPermission(r.Context(), id)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s tried to update permission '%s' but an error happened: %s", user.ID, id, err.Error()), r)
				h.ResNotFound(w, "permission")
				return
			}
			updatedPermission, err := s.Db.Queries.UpdatePermission(r.Context(), sqlc.UpdatePermissionParams{
				ID:               id,
				SecretKeyPattern: dto.SecretKeyPattern,
			})
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to update permission %s: %s", user.ID, id, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(UpdatePermissionEvent, fmt.Sprintf("user %s from %s to %s", user.ID, permission.SecretKeyPattern, updatedPermission.SecretKeyPattern), r)
			}
			h.ResSuccess(w, updatedPermission)
		})

		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			id := chi.URLParam(r, "id")
			_, err := s.Db.Queries.GetPermission(r.Context(), id)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to delete unexisting permission %s: %s", user.ID, id, err.Error()), r)
				h.ResNotFound(w, "permission")
				return
			}

			err = s.Db.Queries.DeletePermission(r.Context(), id)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to delete permission %s: %s", user.ID, id, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(DeleteEvent, fmt.Sprintf("user %s deleted permission %s", user.ID, id), r)
			}
			h.ResSuccess(w, nil)
		})
	})
}
