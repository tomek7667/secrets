package main

import (
	"log/slog"

	"github.com/tomek7667/go-multi-logger-slog/logger"
	"github.com/tomek7667/secrets/secretssdk"
)

var secretsClient *secretssdk.Client

func init() {
	logger.SetLogLevel()

	var err error
	secretsClient, err = secretssdk.New("http://127.0.0.1:7770", "s3cr3t-t0k3n")
	if err != nil {
		panic(err)
	}
}

func main() {
	secret, err := secretsClient.GetSecret("test/key/secret")
	slog.Info(
		"result",
		"secret", secret,
		"err", err,
	)
	allSecrets, err := secretsClient.ListSecrets()
	slog.Info(
		"result",
		"allSecrets", allSecrets,
		"err", err,
	)
}
