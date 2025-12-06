package secrets

import (
	"github.com/go-chi/chi"
	"github.com/tomek7667/go-http-helpers/chii"
)

func (s *Server) AddTokensRoutes() {
	auth := s.Router.With(chii.WithAuth(s.auther))
	auth.Route("/api/tokens", func(r chi.Router) {
		// todo:
	})
}
