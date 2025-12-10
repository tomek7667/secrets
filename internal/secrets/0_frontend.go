package secrets

import (
	"net/http"

	_ "embed"

	"github.com/go-chi/chi"
)

//go:embed index.html
var indexHtml []byte

func (s *Server) AddFrontendRoutes() {
	s.Router.Route("/", func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write(indexHtml)
		})
	})
}
