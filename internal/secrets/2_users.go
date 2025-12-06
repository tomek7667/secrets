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

type CreateUserDto struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UpdateUserDto struct {
	Password string `json:"password"`
}

func (s *Server) AddUsersRoutes() {
	auth := s.Router.With(chii.WithAuth(s.auther))
	auth.Route("/api/users", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			users, err := s.Db.Queries.ListUsers(r.Context())
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("failed to list users for user %s: %s", user.ID, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(GetUsersEvent, fmt.Sprintf("%s retrieved users", user.ID), r)
			}
			h.ResSuccess(w, users)
		})

		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			dto, err := h.GetDto[CreateUserDto](r)
			if err != nil {
				h.ResBadRequest(w, err)
				return
			}
			newuser, err := s.Db.Queries.CreateUser(r.Context(), sqlc.CreateUserParams{
				ID:       utils.CreateUUID(),
				Username: dto.Username,
				Password: dto.Password,
			})
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s tried to create user %s: %s", user.ID, dto.Username, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(IngestEvent, fmt.Sprintf("user %s created user %s %s", user.ID, newuser.ID, newuser.Username), r)
			}
			h.ResSuccess(w, newuser)
		})

		r.Put("/{id}", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			id := chi.URLParam(r, "id")
			dto, err := h.GetDto[UpdateUserDto](r)
			if err != nil {
				h.ResBadRequest(w, err)
				return
			}
			toBeUpdated, err := s.Db.Queries.UpdateUser(r.Context(), sqlc.UpdateUserParams{
				ID:       id,
				Password: dto.Password,
			})
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to update user %s: %s", user.ID, id, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(IngestEvent, fmt.Sprintf("user %s updated user %s", user.ID, id), r)
			}
			h.ResSuccess(w, toBeUpdated)
		})

		r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			id := chi.URLParam(r, "id")
			_, err := s.Db.Queries.GetUser(r.Context(), id)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to delete unexisting user %s: %s", user.ID, id, err.Error()), r)
				h.ResNotFound(w, "user")
				return
			}

			err = s.Db.Queries.DeleteUser(r.Context(), id)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to delete user %s: %s", user.ID, id, err.Error()), r)
				h.ResErr(w, err)
				return
			} else {
				s.Log(DeleteEvent, fmt.Sprintf("user %s deleted user %s", user.ID, id), r)
			}
			h.ResSuccess(w, nil)
		})
	})
}
