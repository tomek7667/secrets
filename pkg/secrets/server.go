package secrets

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/caarlos0/env/v10"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/tomek7667/go-http-helpers/chii"
	"github.com/tomek7667/go-http-helpers/h"
	"github.com/tomek7667/secrets/internal/sqlite"
)

type Server struct {
	Db *sqlite.Client

	Address        string `env:"SECRETS_ADDRESS" envDefault:"127.0.0.1:7770"`
	AllowedOrigins string `env:"ALLOWED_ORIGINS"`
	allowedOrigins []string
	Router         chi.Router
}

func New() (*Server, error) {
	server := &Server{}
	if err := env.Parse(server); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	// db
	godotenv.Load()
	c, err := sqlite.New(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sqlite: %w", err)
	}

	// http
	r := chi.NewRouter()
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		h.ResNotFound(w, "page")
	})

	// init
	server.Db = c
	server.Router = r
	server.allowedOrigins = strings.Split(server.AllowedOrigins, ",")

	return server, nil
}

func (s *Server) Serve() {
	chii.SetupMiddlewares(s.Router, s.allowedOrigins)
	s.SetupRoutes()
	fmt.Printf("listening on %s\n", s.Address)
	chii.PrintRoutes(s.Router)
	err := http.ListenAndServe(s.Address, s.Router)
	if err != nil {
		panic(fmt.Errorf("listen and serve failed s.Address='%s': %w", s.Address, err))
	}
}
