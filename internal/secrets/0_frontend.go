package secrets

import (
	"log/slog"
	"net/http"
	"strings"

	_ "embed"

	"github.com/go-chi/chi"
)

//go:embed index.html
var indexHtml []byte

//go:embed index_withcaptcha.html
var turnstileIndexHtml []byte

func (s *Server) AddFrontendRoutes() {
	htmlToRender := indexHtml
	if s.turnstileSiteKey != "" {
		htmlToRender = []byte(strings.ReplaceAll(string(turnstileIndexHtml), "TURNSTILE_SITE_KEY", s.turnstileSiteKey))
	} else {
		slog.Warn("Rendering login page without captcha: turnstile credentials not configured. Set via -turnstile-secret and -turnstile-site-key flags or TURNSTILE_SECRET/TURNSTILE_SITE_KEY environment variables.")
	}
	s.Router.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(htmlToRender)
		})
	})
}
