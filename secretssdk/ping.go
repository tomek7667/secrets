package secretssdk

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type pingResponse struct {
	Ts string `json:"ts"`
}

func (c *Client) Ping() error {
	endpoint := fmt.Sprintf("%s/ping", c.BaseUrl)

	resp, err := http.Get(endpoint)
	if err != nil {
		return fmt.Errorf("server unreachable: %w", err)
	}
	defer resp.Body.Close()

	var result pingResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("invalid server response: %w", err)
	}

	if result.Ts == "" {
		return fmt.Errorf("server did not return timestamp")
	}

	return nil
}
