package main

import (
	"github.com/tomek7667/secrets/pkg/secrets"
)

func main() {
	server, err := secrets.New()
	if err != nil {
		panic(err)
	}
	server.Serve()
}
