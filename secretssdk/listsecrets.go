package secretssdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tomek7667/go-http-helpers/utils"
)

type secretsResponse struct {
	Success bool     `json:"success"`
	Data    []Secret `json:"data"`
}

func (c *Client) ListSecretsWithCtx(ctx context.Context) (map[string]string, error) {
	endpoint := fmt.Sprintf("%s/api/secrets/list", c.BaseUrl)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new request for endpoint '%s': %w", endpoint, err)
	}
	req = req.WithContext(ctx)
	resp, err := c.GetHttpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed for listing secrets: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: token lacks permission for listing secrets")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for listing secrets", resp.StatusCode)
	}

	var result secretsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response for listing secrets: %w", err)
	}

	r := map[string]string{}
	for _, s := range result.Data {
		val, err := utils.B64Decode(s.Value)
		if err != nil {
			return nil, fmt.Errorf("failed to base64 decode the value for secret key %s: %w", s.Key, err)
		}
		r[s.Key] = val
	}
	return r, nil
}

func (c *Client) ListSecrets() (map[string]string, error) {
	return c.ListSecretsWithCtx(context.Background())
}
