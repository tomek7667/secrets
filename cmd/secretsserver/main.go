// cmd/secretsserver/main.go
package main

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"github.com/tomek7667/secrets/pkg/secrets"
)

type CliOptions struct {
	Address        string `env:"SECRETS_ADDRESS" envDefault:"127.0.0.1:7770"`
	DbPath         string `env:"SECRETS_DB_PATH" envDefault:"./secrets.sqlite"`
	AllowedOrigins string `env:"ALLOWED_ORIGINS"`
}

func main() {
	godotenv.Load()
	var opts CliOptions

	// load defaults from env
	if err := env.Parse(&opts); err != nil {
		log.Fatalf("failed to parse env: %v", err)
	}

	rootCmd := &cobra.Command{
		Use:   "secretsserver",
		Short: "Run secrets HTTP API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			// you’ll make server/client infer from CliOptions
			srv, err := secrets.New(
				opts.Address,
				opts.AllowedOrigins,
				opts.DbPath,
			)
			if err != nil {
				return err
			}
			srv.Serve()
			return nil
		},
	}

	// flags override env/defaults
	rootCmd.Flags().StringVar(&opts.Address, "address", opts.Address, "listen address")
	rootCmd.Flags().StringVar(&opts.DbPath, "db-path", opts.DbPath, "path to sqlite db")
	rootCmd.Flags().StringVar(&opts.AllowedOrigins, "allowed-origins", opts.AllowedOrigins, "comma-separated list of allowed CORS origins")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
