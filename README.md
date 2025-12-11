[![Bruno Integration Tests](https://github.com/tomek7667/secrets/actions/workflows/bruno-tests.yml/badge.svg)](https://github.com/tomek7667/secrets/actions/workflows/bruno-tests.yml)

# Secrets Manager

Self-hosted secrets management with web UI and REST API.

## Features

- 🔐 Web UI with dark mode
- 💾 SQLite storage (single binary, no dependencies)
- 👥 Multi-user with JWT authentication
- 🔑 API tokens with pattern-based permissions
- 📜 Audit logging
- 🔏 Certificate & Key Management (RSA, ECDSA, ED25519)
  - Generate key pairs
  - Import/Export certificates (PEM format)
  - Create self-signed & CA-signed certificates
  - Certificate verification & validation

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
go generate ./...          # Builds frontend + runs migrations
go build ./cmd/secretsserver
```

> **Note:** Building requires Node.js and Yarn for the frontend.

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
| GET/POST/PUT/DELETE | `/api/users`        | Manage users           |
| GET/POST/PUT/DELETE | `/api/tokens`       | Manage tokens          |
| GET/POST/PUT/DELETE | `/api/permissions`  | Manage permissions     |
| GET/POST/DELETE     | `/api/certificates` | Manage certificates    |

## Certificate Management

The secrets manager includes comprehensive certificate and key management capabilities:

### Generate Key Pairs

Generate RSA, ECDSA, or ED25519 key pairs:

```bash
POST /api/certificates/generate-keypair
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "name": "my-keypair",
  "algorithm": "RSA",     // "RSA", "ECDSA", or "ED25519"
  "key_size": 2048        // RSA: 2048, 3072, 4096; ECDSA: 256, 384, 521
}
```

### Generate Certificates

Create self-signed or CA-signed X.509 certificates:

```bash
POST /api/certificates/generate-certificate
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "name": "my-cert",
  "private_key_name": "my-keypair-private",
  "subject": {
    "common_name": "example.com",
    "organization": "My Organization",
    "country": "US"
  },
  "validity_days": 365,
  "is_ca": false,
  "dns_names": ["example.com", "www.example.com"],
  "signing_cert_name": "ca-cert"  // Optional: for CA-signed certs
}
```

### Import/Export Certificates

```bash
# Import a certificate
POST /api/certificates/import
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "name": "imported-cert",
  "cert_type": "certificate",  // "private_key", "public_key", "certificate", "ca_certificate"
  "pem_data": "-----BEGIN CERTIFICATE-----\n..."
}

# Export a certificate
GET /api/certificates/{name}/export
Authorization: Bearer <jwt>
```

### Verify Certificates

Verify certificate validity and signatures:

```bash
POST /api/certificates/verify
Authorization: Bearer <jwt>
Content-Type: application/json

{
  "certificate_name": "my-cert",
  "ca_cert_name": "ca-cert"  // Optional: verify against specific CA
}
```

### Go SDK - Certificate Management

```go
// Generate RSA key pair
keyPair, err := client.GenerateKeyPair(secretssdk.GenerateKeyPairRequest{
    Name:      "my-rsa-key",
    Algorithm: "RSA",
    KeySize:   ptrInt(2048),
})

// Generate self-signed certificate
cert, err := client.GenerateCertificate(secretssdk.GenerateCertificateRequest{
    Name:           "my-cert",
    PrivateKeyName: "my-rsa-key-private",
    Subject: secretssdk.Subject{
        CommonName:   "example.com",
        Organization: "My Org",
        Country:      "US",
    },
    ValidityDays: 365,
    IsCA:         false,
    DNSNames:     []string{"example.com", "*.example.com"},
})

// Verify certificate
result, err := client.VerifyCertificate(secretssdk.VerifyCertificateRequest{
    CertificateName: "my-cert",
})

// Export certificate
exported, err := client.ExportCertificate("my-cert")
fmt.Println(exported.PemData)
```

### List & Manage Certificates

```bash
# List all certificates
GET /api/certificates
Authorization: Bearer <jwt>

# Get specific certificate
GET /api/certificates/{name}
Authorization: Bearer <jwt>

# Delete certificate
DELETE /api/certificates/{name}
Authorization: Bearer <jwt>
```

## Pattern Matching

Permissions use wildcard patterns:

- `*` — all secrets
- `aws/*` — secrets starting with `aws/`
- `exact-key` — exact match only

## Development

```bash
go mod download
cd web && yarn install && yarn build && cd ..
go run cmd/secretsserver/main.go --admin-password "dev123"
```

Frontend development (with hot reload):

```bash
cd web && yarn dev    # Runs on localhost:5173, proxies API to :7770
```

Integration tests (requires [Bruno CLI](https://www.usebruno.com/)):

```bash
# Run all test collections
cd bruno/auth_integration_tests && bru run --env local
cd bruno/secrets_integration_tests && bru run --env local
cd bruno/users_tokens_integration_tests && bru run --env local
cd bruno/certificates_integration_tests && bru run --env local
```
