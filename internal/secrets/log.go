package secrets

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/tomek7667/go-http-helpers/utils"
	"github.com/tomek7667/secrets/internal/sqlc"
)

type LogEvent string

const (
	ErrorEvent            LogEvent = "error"
	UnauthorizedEvent     LogEvent = "unauthorized"
	IngestEvent           LogEvent = "ingest"
	DeleteEvent           LogEvent = "delete"
	GetSecretEvent        LogEvent = "get-secret"
	GetFullEnvEvent       LogEvent = "get-full-env"
	UpdateSecretEvent     LogEvent = "update-secret"
	UpdateTokenEvent      LogEvent = "update-token"
	GetUsersEvent         LogEvent = "get-users"
	GetSecretsEvent       LogEvent = "get-secrets"
	GetTokensEvent        LogEvent = "get-tokens"
	GetPermissionsEvent   LogEvent = "get-permissions"
	UpdatePermissionEvent LogEvent = "update-permission"
	LoginSuccessEvent     LogEvent = "login-success"
	LoginFailedEvent      LogEvent = "login-failed"
)

func (le LogEvent) String() string {
	return string(le)
}

func (s *Server) Log(event LogEvent, msg string, r *http.Request) {
	go func() {
		requestedUrl := r.Method + " " + r.URL.String()
		slog.Debug(
			msg,
			"event", event,
		)

		_, err := s.Db.Queries.CreateLog(context.Background(), sqlc.CreateLogParams{
			ID:           utils.CreateUUID(),
			Event:        event.String(),
			Msg:          msg,
			RequestedUrl: &requestedUrl,
			RemoteAddr:   &r.RemoteAddr,
		})
		if err != nil {
			slog.Error(
				"failed to save a log entry",
				"err", err,
				"event", event,
				"msg", msg,
			)
		}
	}()
}
