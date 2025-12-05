package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tomek7667/go-http-helpers/utils"
	"github.com/tomek7667/secrets/internal/sqlc"
	"github.com/tomek7667/secrets/internal/sqlite"
)

type Auther struct {
	Db        *sqlite.Client
	JwtSecret string
}

func (a Auther) GetUserFromToken(ctx context.Context, token string) (*sqlc.User, error) {
	user, err := a.Db.Queries.GetUser(ctx, token)
	if err != nil {
		return nil, fmt.Errorf("get user failed for token '%s': %w", token, err)
	}
	return &user, nil
}

func (a Auther) GetToken(user *sqlc.User) (string, error) {
	var umap map[string]any
	b, err := json.Marshal(user)
	if err != nil {
		return "", fmt.Errorf("failed to json.marshal the user: %w", err)
	}

	err = json.Unmarshal(b, &umap)
	if err != nil {
		return "", fmt.Errorf("failed to json.unmarshal the user: %w", err)
	}
	umap["created"] = time.Now()

	return utils.JwtEncode(umap, a.JwtSecret)
}
