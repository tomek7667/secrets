package secretssdk

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Secret struct {
	ID    string `json:"id"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

type secretResponse struct {
	Success bool   `json:"success"`
	Data    Secret `json:"data"`
}

func (c *Client) GetSecretWithCtx(key string, ctx context.Context) (*Secret, error) {
	endpoint := fmt.Sprintf("%s/api/secrets/get?key=%s", c.BaseUrl, url.QueryEscape(key))
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create a new request for endpoint '%s': %w", endpoint, err)
	}
	req = req.WithContext(ctx)
	resp, err := c.GetHttpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed for secret '%s': %w", key, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("unauthorized: token lacks permission for secret '%s'", key)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d for secret '%s'", resp.StatusCode, key)
	}

	var result secretResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response for secret '%s': %w", key, err)
	}

	decoded, err := base64.StdEncoding.DecodeString(result.Data.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret value: %w", err)
	}
	result.Data.Value = string(decoded)

	return &result.Data, nil
}

func (c *Client) GetSecret(key string) (*Secret, error) {
	return c.GetSecretWithCtx(key, context.Background())
}
