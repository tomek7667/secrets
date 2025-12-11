package secrets

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
)

type turnstileResp struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes,omitempty"`
}

func (s *Server) verifyCaptcha(ctx context.Context, token string) error {
	if s.turnstileSecret == "" {
		slog.Info("Skipping captcha verification: turnstile secret not configured. Set via -turnstile-secret flag or TURNSTILE_SECRET environment variable.")
		return nil
	}
	if token == "" {
		return fmt.Errorf("missing captcha")
	}

	form := url.Values{}
	form.Set("secret", s.turnstileSecret)
	form.Set("response", token)

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		"https://challenges.cloudflare.com/turnstile/v0/siteverify",
		strings.NewReader(form.Encode()),
	)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code from captcha verification: %d", resp.StatusCode)
	}

	var body turnstileResp
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return fmt.Errorf("failed to decode captcha verification response: %w", err)
	}
	if !body.Success {
		return fmt.Errorf("invalid captcha: %v", body.ErrorCodes)
	}
	return nil
}
