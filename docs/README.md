# Documentation

This directory contains documentation assets for the Secrets Manager project.

## Structure

```
docs/
├── README.md           # This file
└── screenshots/        # Application screenshots for documentation
    ├── 01-login.png              # Login page
    ├── 02-secrets-dashboard.png  # Secrets management panel
    ├── 03-secret-revealed.png    # Secret value revealed
    ├── 04-users-panel.png        # User management panel
    ├── 05-tokens-panel.png       # API token management panel
    ├── 06-permissions-panel.png  # Permission configuration panel
    └── 07-login-with-captcha.png # Login with Cloudflare Turnstile captcha
```

## Screenshots

All screenshots are taken at 1280x720 resolution in PNG format and show the application's dark mode interface.

### Adding New Screenshots

When updating screenshots:

1. Start the server: `go run cmd/secretsserver/main.go --admin-password "admin123"`
2. Navigate to `http://127.0.0.1:7770`
3. Capture screenshots of each panel
4. Save them with the naming convention: `##-descriptive-name.png`
5. Update the main README.md to reference new screenshots

### Screenshot Guidelines

- Use consistent viewport size (1280x720 or similar)
- Capture the full application interface including header and navigation
- Use realistic but safe example data (no real secrets)
- Ensure dark mode is enabled for consistency
