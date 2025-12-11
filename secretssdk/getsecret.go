package secretssdk

import (
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

func (c *Client) GetSecret(key string) (*Secret, error) {
	endpoint := fmt.Sprintf("%s/api/secrets/get?key=%s", c.BaseUrl, url.QueryEscape(key))

	resp, err := c.GetHttpClient().Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get secret: status %d", resp.StatusCode)
	}

	var result secretResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	decoded, err := base64.StdEncoding.DecodeString(result.Data.Value)
	if err != nil {
		return nil, fmt.Errorf("failed to decode secret value: %w", err)
	}
	result.Data.Value = string(decoded)

	return &result.Data, nil
}
