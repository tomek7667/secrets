# Secrets Manager

Self-hosted secrets management with web UI and REST API.

## Features

- Web UI with dark mode
- SQLite storage (single binary, no dependencies)
- Multi-user with JWT authentication
- API tokens with pattern-based permissions
- Audit logging

## Screenshots

![Secrets Dashboard](docs/screenshots/02-secrets-dashboard.png)
![Permissions](docs/screenshots/06-permissions.png)

## Installation

```bash
go install github.com/tomek7667/secrets/cmd/secretsserver@latest
```

Or build from source:

```bash
git clone https://github.com/tomek7667/secrets.git
cd secrets
go build ./cmd/secretsserver
```

## Usage

```bash
secretsserver
```

Default: `http://127.0.0.1:7770`. Admin credentials are logged on first run.

### Configuration

| Flag / Env                                    | Default            | Description            |
| --------------------------------------------- | ------------------ | ---------------------- |
| `--address` / `SECRETS_ADDRESS`               | `127.0.0.1:7770`   | Listen address         |
| `--db-path` / `SECRETS_DB_PATH`               | `./secrets.sqlite` | SQLite database path   |
| `--jwt-secret` / `SECRETS_JWT_SECRET`         | (auto)             | JWT signing secret     |
| `--admin-password` / `SECRETS_ADMIN_PASSWORD` | (auto)             | Initial admin password |
| `--allowed-origins` / `ALLOWED_ORIGINS`       | (none)             | CORS origins           |

## Go SDK

```bash
go get github.com/tomek7667/secrets/secretssdk
```

```go
package main

import (
    "fmt"
    "log"

    "github.com/tomek7667/secrets/secretssdk"
)

func main() {
    client, err := secretssdk.New("http://127.0.0.1:7770", "your-api-token")
    if err != nil {
        log.Fatal(err)
    }

    secret, err := client.GetSecret("my-secret-key")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(secret.Value)
}
```

## API

### Get Secret (API Token)

```bash
GET /api/secrets/get?key=my-secret
Authorization: Api <token>
```

### JWT-Protected Endpoints

Login first:

```bash
POST /login
{"username": "admin", "password": "..."}
```

Then use `Authorization: Bearer <jwt>` for:

| Method              | Endpoint            | Description        |
| ------------------- | ------------------- | ------------------ |
| GET                 | `/api/secrets`      | List secrets       |
| POST                | `/api/secrets`      | Create secret      |
| PUT                 | `/api/secrets?key=` | Update secret      |
| DELETE              | `/api/secrets?key=` | Delete secret      |
| GET/POST/PUT/DELETE | `/api/users`        | Manage users       |
| GET/POST/PUT/DELETE | `/api/tokens`       | Manage tokens      |
| GET/POST/PUT/DELETE | `/api/permissions`  | Manage permissions |

## Pattern Matching

Permissions use wildcard patterns:

- `*` — all secrets
- `aws/*` — secrets starting with `aws/`
- `exact-key` — exact match only

## Development

```bash
go mod download
go run cmd/secretsserver/main.go --admin-password "dev123"
```

Integration tests (requires [Bruno CLI](https://www.usebruno.com/)):

```bash
cd bruno/auth_integration_tests && bru run --env local
```
