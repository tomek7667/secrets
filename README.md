# secrets

secrets manager ; store passwords, api keys, files locally with a backup online

```bash
go install github.com/tomek7667/secrets/cmd/secretsserver@latest
```

---

## envs

- `SECRETS_DB_PATH` (`./secrets.sqlite`)
- `SECRETS_FRONTEND_BASE_URL` (`http://127.0.0.1:7770`)
- `ALLOWED_ORIGINS` (no default) - comma separated CORS origins

## development

build:

`go build ./cmd/secretsserver`
