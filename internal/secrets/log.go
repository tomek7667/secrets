package secrets

import (
	"context"
	"log/slog"

	"github.com/tomek7667/go-http-helpers/utils"
	"github.com/tomek7667/secrets/internal/sqlc"
)

type LogEvent string

const (
	ErrorEvent     LogEvent = "error"
	IngestEvent    LogEvent = "ingest"
	GetSecretEvent LogEvent = "get-secret"
)

func (le LogEvent) String() string {
	return string(le)
}

func (s *Server) Log(event LogEvent, msg string) {
	go func() {
		slog.Debug(
			msg,
			"event", event,
		)
		_, err := s.Db.Queries.CreateLog(context.Background(), sqlc.CreateLogParams{
			ID:    utils.CreateUUID(),
			Event: event.String(),
			Msg:   msg,
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
