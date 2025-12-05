package secrets

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/tomek7667/go-http-helpers/chii"
	"github.com/tomek7667/go-http-helpers/h"
	"github.com/tomek7667/secrets/internal/sqlite"
)

type Options struct {
	Address        string
	DBPath         string
	AllowedOrigins string
}

type Server struct {
	Db *sqlite.Client

	Address        string
	allowedOrigins []string
	Router         chi.Router
}

func New(address, allowedOrigins, dbPath string) (*Server, error) {
	// db
	godotenv.Load()
	c, err := sqlite.New(context.Background(), dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize sqlite: %w", err)
	}

	// http
	r := chi.NewRouter()
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		h.ResNotFound(w, "page")
	})

	// init
	server := &Server{
		Address:        address,
		allowedOrigins: strings.Split(allowedOrigins, ","),
		Db:             c,
		Router:         r,
	}

	return server, nil
}

func (s *Server) Serve() {
	chii.SetupMiddlewares(s.Router, s.allowedOrigins)
	s.SetupRoutes()
	fmt.Printf("listening on address '%s'\n", s.Address)
	chii.PrintRoutes(s.Router)
	err := http.ListenAndServe(s.Address, s.Router)
	if err != nil {
		panic(fmt.Errorf("listen and serve failed s.Address='%s': %w", s.Address, err))
	}
}
