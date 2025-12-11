# Certificate Management Fix Guide

This document explains how to fix the critical authentication issue in the certificate management feature.

## Problem

The certificate endpoints currently use JWT authentication middleware:
```go
auth := s.Router.With(chii.WithAuth(s.auther))
auth.Route("/api/certificates", func(r chi.Router) {
    // All routes here require JWT auth
})
```

This means:
- SDK cannot use API tokens to access certificates
- Only JWT (username/password login) works  
- Defeats the purpose of token-based permissions

## Solution

Follow the pattern from `/api/secrets/get` (lines 117-162 in `1_secrets.go`):

1. Remove JWT middleware wrapper
2. Each handler validates API token manually
3. Check permissions against certificate name patterns

## Example: Fix GET /api/certificates/{name}

### Before (BROKEN):
```go
func (s *Server) AddCertificatesRoutes() {
    auth := s.Router.With(chii.WithAuth(s.auther)) // JWT wrapper
    auth.Route("/api/certificates", func(r chi.Router) {
        r.Get("/{name}", func(w http.ResponseWriter, r *http.Request) {
            user := chii.GetUser[sqlc.User](r) // Assumes JWT
            name := chi.URLParam(r, "name")
            certificate, err := s.Db.Queries.GetCertificate(r.Context(), name)
            // ... rest of handler
        })
    })
}
```

### After (FIXED):
```go
func (s *Server) AddCertificatesRoutes() {
    // Remove JWT wrapper - use s.Router directly
    s.Router.Route("/api/certificates", func(r chi.Router) {
        r.Get("/{name}", func(w http.ResponseWriter, r *http.Request) {
            // 1. Parse and validate API token
            authValue := strings.TrimSpace(r.Header.Get("Authorization"))
            name := chi.URLParam(r, "name")
            
            if !strings.HasPrefix(authValue, "Api ") {
                s.Log(UnauthorizedEvent, fmt.Sprintf("invalid token type for certificate '%s'", name), r)
                h.ResUnauthorized(w)
                return
            }
            
            token, _ := strings.CutPrefix(authValue, "Api ")
            tkn, err := s.Db.Queries.GetTokenByToken(r.Context(), token)
            if err != nil {
                s.Log(UnauthorizedEvent, fmt.Sprintf("invalid token for certificate '%s': %s", name, err.Error()), r)
                h.ResUnauthorized(w)
                return
            }
            
            // 2. Get certificate (before permission check to know what we're checking)
            certificate, err := s.Db.Queries.GetCertificate(r.Context(), name)
            if err != nil {
                s.Log(ErrorEvent, fmt.Sprintf("failed to get certificate '%s': %s", name, err.Error()), r)
                h.ResNotFound(w, "certificate")
                return
            }
            
            // 3. Check permissions
            permissions, err := s.Db.Queries.ListPermissionsByTokenId(r.Context(), tkn.ID)
            if err != nil {
                s.Log(ErrorEvent, fmt.Sprintf("failed to list permissions for token %s: %s", tkn.ID, err.Error()), r)
                h.ResErr(w, err)
                return
            }
            
            matched := false
            certPattern := fmt.Sprintf("cert:%s", certificate.Name) // Use "cert:" prefix
            for _, permission := range permissions {
                if PatternMatches(certPattern, permission.SecretKeyPattern) {
                    matched = true
                    break
                }
            }
            
            if !matched {
                s.Log(UnauthorizedEvent, fmt.Sprintf("token %s cannot access certificate %s", tkn.ID, name), r)
                h.ResUnauthorized(w)
                return
            }
            
            // 4. Return the certificate
            s.Log(GetSecretsEvent, fmt.Sprintf("token %s retrieved certificate %s", tkn.ID, name), r)
            h.ResSuccess(w, certificate)
        })
    })
}
```

## Permission Pattern

For certificates, use the `cert:` prefix in permission patterns:
- `cert:*` - access all certificates
- `cert:myapp-*` - access certificates starting with "myapp-"
- `cert:my-specific-cert` - access only specific certificate

This distinguishes certificate permissions from secret permissions.

## Implementation Checklist

For EACH endpoint in `AddCertificatesRoutes()`:

1. [ ] Remove JWT auth wrapper (line 74)
2. [ ] Add token parsing boilerplate (lines from example above)
3. [ ] Get the resource BEFORE permission check (to know what pattern to match)
4. [ ] Check permissions using `PatternMatches()` with `cert:` prefix
5. [ ] Update log statements to use `tkn.ID` instead of `user.ID`
6. [ ] Test with Bruno tests

## Endpoints to Fix

- [  ] GET `/` (list all certificates)
- [ ] GET `/{name}` (get specific certificate)
- [ ] POST `/generate-keypair` (generate key pair)
- [ ] POST `/import` (import certificate)
- [ ] GET `/{name}/export` (export certificate)
- [ ] POST `/generate-certificate` (generate certificate)
- [ ] POST `/verify` (verify certificate)
- [ ] DELETE `/{name}` (delete certificate)

## Testing

After fixing, run Bruno tests:
```bash
cd bruno/certificates_integration_tests
bru run --env local
```

All 18 tests should pass with API token authentication.

## Notes

- Private keys are named with `-private` suffix automatically
- Certificate names should be used in permission patterns with `cert:` prefix
- This matches the existing `/api/secrets/get` pattern for consistency
