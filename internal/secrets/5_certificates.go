package secrets

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/tomek7667/go-http-helpers/chii"
	"github.com/tomek7667/go-http-helpers/h"
	"github.com/tomek7667/go-http-helpers/utils"
	"github.com/tomek7667/secrets/internal/sqlc"
)

type GenerateKeyPairDto struct {
	Name      string `json:"name"`
	Algorithm string `json:"algorithm"` // "RSA", "ECDSA", "ED25519"
	KeySize   *int   `json:"key_size"`  // For RSA: 2048, 3072, 4096; For ECDSA: 256, 384, 521
}

type ImportCertificateDto struct {
	Name     string `json:"name"`
	CertType string `json:"cert_type"` // "private_key", "public_key", "certificate", "ca_certificate"
	PemData  string `json:"pem_data"`
}

type GenerateCertificateDto struct {
	Name            string   `json:"name"`
	PrivateKeyName  string   `json:"private_key_name"`
	Subject         Subject  `json:"subject"`
	ValidityDays    int      `json:"validity_days"`
	IsCA            bool     `json:"is_ca"`
	DNSNames        []string `json:"dns_names,omitempty"`
	EmailAddresses  []string `json:"email_addresses,omitempty"`
	SigningCertName *string  `json:"signing_cert_name,omitempty"` // For CA-signed certificates
}

type Subject struct {
	CommonName         string `json:"common_name"`
	Organization       string `json:"organization,omitempty"`
	OrganizationalUnit string `json:"organizational_unit,omitempty"`
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	Locality           string `json:"locality,omitempty"`
}

type VerifyCertificateDto struct {
	CertificateName string  `json:"certificate_name"`
	CACertName      *string `json:"ca_cert_name,omitempty"` // Optional: verify against specific CA
}

type CertificateMetadata struct {
	Issuer       string    `json:"issuer,omitempty"`
	Subject      string    `json:"subject,omitempty"`
	NotBefore    time.Time `json:"not_before,omitempty"`
	NotAfter     time.Time `json:"not_after,omitempty"`
	SerialNumber string    `json:"serial_number,omitempty"`
	IsCA         bool      `json:"is_ca,omitempty"`
	DNSNames     []string  `json:"dns_names,omitempty"`
	KeyUsage     []string  `json:"key_usage,omitempty"`
}

func (s *Server) AddCertificatesRoutes() {
	auth := s.Router.With(chii.WithAuth(s.auther))
	auth.Route("/api/certificates", func(r chi.Router) {
		// List all certificates
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			certificates, err := s.Db.Queries.ListCertificates(r.Context())
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("failed to list certificates for user %s: %s", user.ID, err.Error()), r)
				h.ResErr(w, err)
				return
			}
			s.Log(GetSecretsEvent, fmt.Sprintf("%s retrieved certificates", user.ID), r)
			h.ResSuccess(w, certificates)
		})

		// Get certificate by name
		r.Get("/{name}", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			name := chi.URLParam(r, "name")
			certificate, err := s.Db.Queries.GetCertificate(r.Context(), name)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s tried to get certificate '%s' but an error happened: %s", user.ID, name, err.Error()), r)
				h.ResNotFound(w, "certificate")
				return
			}
			s.Log(GetSecretsEvent, fmt.Sprintf("%s retrieved certificate %s", user.ID, name), r)
			h.ResSuccess(w, certificate)
		})

		// Generate key pair
		r.Post("/generate-keypair", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			dto, err := h.GetDto[GenerateKeyPairDto](r)
			if err != nil {
				h.ResBadRequest(w, err)
				return
			}

			privateKeyPEM, publicKeyPEM, algorithm, keySize, err := generateKeyPair(dto.Algorithm, dto.KeySize)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to generate key pair: %s", user.ID, err.Error()), r)
				h.ResErr(w, err)
				return
			}

			// Store private key
			privateKeyCert, err := s.Db.Queries.CreateCertificate(r.Context(), sqlc.CreateCertificateParams{
				ID:        utils.CreateUUID(),
				Name:      dto.Name + "-private",
				CertType:  "private_key",
				Algorithm: algorithm,
				KeySize:   &keySize,
				PemData:   privateKeyPEM,
				Metadata:  nil,
			})
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to store private key: %s", user.ID, err.Error()), r)
				h.ResErr(w, err)
				return
			}

			// Store public key
			publicKeyCert, err := s.Db.Queries.CreateCertificate(r.Context(), sqlc.CreateCertificateParams{
				ID:        utils.CreateUUID(),
				Name:      dto.Name + "-public",
				CertType:  "public_key",
				Algorithm: algorithm,
				KeySize:   &keySize,
				PemData:   publicKeyPEM,
				Metadata:  nil,
			})
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to store public key: %s", user.ID, err.Error()), r)
				h.ResErr(w, err)
				return
			}

			s.Log(IngestEvent, fmt.Sprintf("user %s generated key pair %s", user.ID, dto.Name), r)
			h.ResSuccess(w, map[string]interface{}{
				"private_key": privateKeyCert,
				"public_key":  publicKeyCert,
			})
		})

		// Import certificate
		r.Post("/import", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			dto, err := h.GetDto[ImportCertificateDto](r)
			if err != nil {
				h.ResBadRequest(w, err)
				return
			}

			// Parse and validate PEM data
			block, _ := pem.Decode([]byte(dto.PemData))
			if block == nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s tried to import invalid PEM data", user.ID), r)
				h.ResBadRequest(w, fmt.Errorf("invalid PEM data"))
				return
			}

			algorithm, keySize, metadata, err := parsePEMBlock(block, dto.CertType)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to parse PEM: %s", user.ID, err.Error()), r)
				h.ResErr(w, err)
				return
			}

			metadataJSON, _ := json.Marshal(metadata)
			metadataStr := string(metadataJSON)

			certificate, err := s.Db.Queries.CreateCertificate(r.Context(), sqlc.CreateCertificateParams{
				ID:        utils.CreateUUID(),
				Name:      dto.Name,
				CertType:  dto.CertType,
				Algorithm: algorithm,
				KeySize:   keySize,
				PemData:   dto.PemData,
				Metadata:  &metadataStr,
			})
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to import certificate: %s", user.ID, err.Error()), r)
				h.ResErr(w, err)
				return
			}

			s.Log(IngestEvent, fmt.Sprintf("user %s imported certificate %s", user.ID, dto.Name), r)
			h.ResSuccess(w, certificate)
		})

		// Export certificate (returns PEM data)
		r.Get("/{name}/export", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			name := chi.URLParam(r, "name")
			certificate, err := s.Db.Queries.GetCertificate(r.Context(), name)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s tried to export certificate '%s' but an error happened: %s", user.ID, name, err.Error()), r)
				h.ResNotFound(w, "certificate")
				return
			}
			s.Log(GetSecretsEvent, fmt.Sprintf("%s exported certificate %s", user.ID, name), r)
			h.ResSuccess(w, map[string]string{
				"name":     certificate.Name,
				"pem_data": certificate.PemData,
			})
		})

		// Generate certificate
		r.Post("/generate-certificate", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			dto, err := h.GetDto[GenerateCertificateDto](r)
			if err != nil {
				h.ResBadRequest(w, err)
				return
			}

			// Get private key
			privateKeyCert, err := s.Db.Queries.GetCertificate(r.Context(), dto.PrivateKeyName)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s tried to use non-existent private key '%s'", user.ID, dto.PrivateKeyName), r)
				h.ResNotFound(w, "private_key")
				return
			}

			// Parse private key
			privateKey, err := parsePrivateKey(privateKeyCert.PemData)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to parse private key: %s", user.ID, err.Error()), r)
				h.ResErr(w, err)
				return
			}

			var signingCert *x509.Certificate
			var signingKey interface{}

			if dto.SigningCertName != nil {
				// CA-signed certificate
				signingCertData, err := s.Db.Queries.GetCertificate(r.Context(), *dto.SigningCertName)
				if err != nil {
					s.Log(ErrorEvent, fmt.Sprintf("user %s tried to use non-existent CA certificate '%s'", user.ID, *dto.SigningCertName), r)
					h.ResNotFound(w, "ca_certificate")
					return
				}

				signingCert, err = parseCertificate(signingCertData.PemData)
				if err != nil {
					s.Log(ErrorEvent, fmt.Sprintf("user %s failed to parse CA certificate: %s", user.ID, err.Error()), r)
					h.ResErr(w, err)
					return
				}

				// Get CA private key
				// NOTE: This assumes the CA's private key follows the naming convention: {ca_cert_name}-private
				// For example, if the CA certificate is named "my-ca", its private key must be named "my-ca-private"
				// This is automatically handled when using the generate-keypair endpoint, which creates keys with -private suffix
				caPrivateKeyName := *dto.SigningCertName + "-private"
				caPrivateKeyCert, err := s.Db.Queries.GetCertificate(r.Context(), caPrivateKeyName)
				if err != nil {
					s.Log(ErrorEvent, fmt.Sprintf("user %s tried to use non-existent CA private key '%s'", user.ID, caPrivateKeyName), r)
					h.ResNotFound(w, "ca_private_key")
					return
				}

				signingKey, err = parsePrivateKey(caPrivateKeyCert.PemData)
				if err != nil {
					s.Log(ErrorEvent, fmt.Sprintf("user %s failed to parse CA private key: %s", user.ID, err.Error()), r)
					h.ResErr(w, err)
					return
				}
			}

			certPEM, metadata, err := generateCertificate(*dto, privateKey, signingCert, signingKey)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to generate certificate: %s", user.ID, err.Error()), r)
				h.ResErr(w, err)
				return
			}

			metadataJSON, _ := json.Marshal(metadata)
			metadataStr := string(metadataJSON)

			certType := "certificate"
			if dto.IsCA {
				certType = "ca_certificate"
			}

			var keySize *int64
			if privateKeyCert.KeySize != nil {
				keySize = privateKeyCert.KeySize
			}

			certificate, err := s.Db.Queries.CreateCertificate(r.Context(), sqlc.CreateCertificateParams{
				ID:        utils.CreateUUID(),
				Name:      dto.Name,
				CertType:  certType,
				Algorithm: privateKeyCert.Algorithm,
				KeySize:   keySize,
				PemData:   certPEM,
				Metadata:  &metadataStr,
			})
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to store certificate: %s", user.ID, err.Error()), r)
				h.ResErr(w, err)
				return
			}

			s.Log(IngestEvent, fmt.Sprintf("user %s generated certificate %s", user.ID, dto.Name), r)
			h.ResSuccess(w, certificate)
		})

		// Verify certificate
		r.Post("/verify", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			dto, err := h.GetDto[VerifyCertificateDto](r)
			if err != nil {
				h.ResBadRequest(w, err)
				return
			}

			// Get certificate to verify
			certData, err := s.Db.Queries.GetCertificate(r.Context(), dto.CertificateName)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s tried to verify non-existent certificate '%s'", user.ID, dto.CertificateName), r)
				h.ResNotFound(w, "certificate")
				return
			}

			cert, err := parseCertificate(certData.PemData)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s failed to parse certificate: %s", user.ID, err.Error()), r)
				h.ResErr(w, err)
				return
			}

			var verificationResult map[string]interface{}

			if dto.CACertName != nil {
				// Verify against specific CA
				caCertData, err := s.Db.Queries.GetCertificate(r.Context(), *dto.CACertName)
				if err != nil {
					s.Log(ErrorEvent, fmt.Sprintf("user %s tried to use non-existent CA certificate '%s'", user.ID, *dto.CACertName), r)
					h.ResNotFound(w, "ca_certificate")
					return
				}

				caCert, err := parseCertificate(caCertData.PemData)
				if err != nil {
					s.Log(ErrorEvent, fmt.Sprintf("user %s failed to parse CA certificate: %s", user.ID, err.Error()), r)
					h.ResErr(w, err)
					return
				}

				verificationResult = verifyCertificate(cert, caCert)
			} else {
				// Self-verification (check if self-signed is valid)
				verificationResult = verifyCertificate(cert, nil)
			}

			s.Log(GetSecretsEvent, fmt.Sprintf("%s verified certificate %s", user.ID, dto.CertificateName), r)
			h.ResSuccess(w, verificationResult)
		})

		// Delete certificate
		r.Delete("/{name}", func(w http.ResponseWriter, r *http.Request) {
			user := chii.GetUser[sqlc.User](r)
			name := chi.URLParam(r, "name")
			err := s.Db.Queries.DeleteCertificate(r.Context(), name)
			if err != nil {
				s.Log(ErrorEvent, fmt.Sprintf("user %s tried to delete certificate '%s' but an error happened: %s", user.ID, name, err.Error()), r)
				h.ResNotFound(w, "certificate")
				return
			}
			s.Log(DeleteEvent, fmt.Sprintf("%s deleted certificate %s", user.ID, name), r)
			h.ResSuccess(w, map[string]string{"message": "certificate deleted"})
		})
	})
}

// Helper functions

func generateKeyPair(algorithm string, keySize *int) (privateKeyPEM, publicKeyPEM, algo string, size int64, err error) {
	switch algorithm {
	case "RSA":
		size := 2048
		if keySize != nil {
			size = *keySize
		}
		privateKey, err := rsa.GenerateKey(rand.Reader, size)
		if err != nil {
			return "", "", "", 0, err
		}
		privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
		privateKeyPEM = string(pem.EncodeToMemory(&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
		}))

		publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		if err != nil {
			return "", "", "", 0, err
		}
		publicKeyPEM = string(pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		}))
		return privateKeyPEM, publicKeyPEM, "RSA", int64(size), nil

	case "ECDSA":
		var curve elliptic.Curve
		size := 256
		if keySize != nil {
			size = *keySize
		}
		switch size {
		case 256:
			curve = elliptic.P256()
		case 384:
			curve = elliptic.P384()
		case 521:
			curve = elliptic.P521()
		default:
			return "", "", "", 0, fmt.Errorf("unsupported ECDSA key size: %d", size)
		}

		privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return "", "", "", 0, err
		}
		privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			return "", "", "", 0, err
		}
		privateKeyPEM = string(pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: privateKeyBytes,
		}))

		publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		if err != nil {
			return "", "", "", 0, err
		}
		publicKeyPEM = string(pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		}))
		return privateKeyPEM, publicKeyPEM, "ECDSA", int64(size), nil

	case "ED25519":
		publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return "", "", "", 0, err
		}
		privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
		if err != nil {
			return "", "", "", 0, err
		}
		privateKeyPEM = string(pem.EncodeToMemory(&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: privateKeyBytes,
		}))

		publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
		if err != nil {
			return "", "", "", 0, err
		}
		publicKeyPEM = string(pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		}))
		return privateKeyPEM, publicKeyPEM, "ED25519", 256, nil

	default:
		return "", "", "", 0, fmt.Errorf("unsupported algorithm: %s", algorithm)
	}
}

func parsePEMBlock(block *pem.Block, certType string) (algorithm string, keySize *int64, metadata *CertificateMetadata, err error) {
	metadata = &CertificateMetadata{}

	switch certType {
	case "certificate", "ca_certificate":
		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return "", nil, nil, err
		}
		metadata.Issuer = cert.Issuer.String()
		metadata.Subject = cert.Subject.String()
		metadata.NotBefore = cert.NotBefore
		metadata.NotAfter = cert.NotAfter
		metadata.SerialNumber = cert.SerialNumber.String()
		metadata.IsCA = cert.IsCA
		metadata.DNSNames = cert.DNSNames

		switch cert.PublicKey.(type) {
		case *rsa.PublicKey:
			algorithm = "RSA"
			rsaKey := cert.PublicKey.(*rsa.PublicKey)
			size := int64(rsaKey.N.BitLen())
			keySize = &size
		case *ecdsa.PublicKey:
			algorithm = "ECDSA"
			ecKey := cert.PublicKey.(*ecdsa.PublicKey)
			size := int64(ecKey.Params().BitSize)
			keySize = &size
		case ed25519.PublicKey:
			algorithm = "ED25519"
			size := int64(256)
			keySize = &size
		}
		return algorithm, keySize, metadata, nil

	case "private_key":
		if block.Type == "RSA PRIVATE KEY" {
			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return "", nil, nil, err
			}
			algorithm = "RSA"
			size := int64(key.N.BitLen())
			keySize = &size
			return algorithm, keySize, nil, nil
		} else if block.Type == "EC PRIVATE KEY" {
			key, err := x509.ParseECPrivateKey(block.Bytes)
			if err != nil {
				return "", nil, nil, err
			}
			algorithm = "ECDSA"
			size := int64(key.Params().BitSize)
			keySize = &size
			return algorithm, keySize, nil, nil
		} else if block.Type == "PRIVATE KEY" {
			// PKCS8 format (ED25519 uses this)
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return "", nil, nil, err
			}
			switch key.(type) {
			case ed25519.PrivateKey:
				algorithm = "ED25519"
				size := int64(256)
				keySize = &size
				return algorithm, keySize, nil, nil
			case *rsa.PrivateKey:
				rsaKey := key.(*rsa.PrivateKey)
				algorithm = "RSA"
				size := int64(rsaKey.N.BitLen())
				keySize = &size
				return algorithm, keySize, nil, nil
			case *ecdsa.PrivateKey:
				ecKey := key.(*ecdsa.PrivateKey)
				algorithm = "ECDSA"
				size := int64(ecKey.Params().BitSize)
				keySize = &size
				return algorithm, keySize, nil, nil
			}
		}

	case "public_key":
		publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return "", nil, nil, err
		}
		switch publicKey.(type) {
		case *rsa.PublicKey:
			algorithm = "RSA"
			rsaKey := publicKey.(*rsa.PublicKey)
			size := int64(rsaKey.N.BitLen())
			keySize = &size
		case *ecdsa.PublicKey:
			algorithm = "ECDSA"
			ecKey := publicKey.(*ecdsa.PublicKey)
			size := int64(ecKey.Params().BitSize)
			keySize = &size
		case ed25519.PublicKey:
			algorithm = "ED25519"
			size := int64(256)
			keySize = &size
		}
		return algorithm, keySize, nil, nil
	}

	return "", nil, nil, fmt.Errorf("unsupported certificate type or PEM block")
}

func parsePrivateKey(pemData string) (interface{}, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	if block.Type == "RSA PRIVATE KEY" {
		return x509.ParsePKCS1PrivateKey(block.Bytes)
	} else if block.Type == "EC PRIVATE KEY" {
		return x509.ParseECPrivateKey(block.Bytes)
	} else if block.Type == "PRIVATE KEY" {
		return x509.ParsePKCS8PrivateKey(block.Bytes)
	}

	return nil, fmt.Errorf("unsupported private key type")
}

func parseCertificate(pemData string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(pemData))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}
	return x509.ParseCertificate(block.Bytes)
}

func generateCertificate(dto GenerateCertificateDto, privateKey interface{}, signingCert *x509.Certificate, signingKey interface{}) (string, *CertificateMetadata, error) {
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if err != nil {
		return "", nil, err
	}

	subject := pkix.Name{
		CommonName:         dto.Subject.CommonName,
		Organization:       []string{dto.Subject.Organization},
		OrganizationalUnit: []string{dto.Subject.OrganizationalUnit},
		Country:            []string{dto.Subject.Country},
		Province:           []string{dto.Subject.Province},
		Locality:           []string{dto.Subject.Locality},
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(time.Duration(dto.ValidityDays) * 24 * time.Hour)

	template := x509.Certificate{
		SerialNumber:          serialNumber,
		Subject:               subject,
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		DNSNames:              dto.DNSNames,
		EmailAddresses:        dto.EmailAddresses,
	}

	if dto.IsCA {
		template.IsCA = true
		template.KeyUsage |= x509.KeyUsageCertSign
	}

	var publicKey interface{}
	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		publicKey = &key.PublicKey
	case *ecdsa.PrivateKey:
		publicKey = &key.PublicKey
	case ed25519.PrivateKey:
		publicKey = key.Public()
	default:
		return "", nil, fmt.Errorf("unsupported private key type")
	}

	var parent *x509.Certificate
	var signingPrivateKey interface{}

	if signingCert != nil && signingKey != nil {
		// CA-signed certificate
		parent = signingCert
		signingPrivateKey = signingKey
	} else {
		// Self-signed certificate
		parent = &template
		signingPrivateKey = privateKey
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, parent, publicKey, signingPrivateKey)
	if err != nil {
		return "", nil, err
	}

	certPEM := string(pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	}))

	metadata := &CertificateMetadata{
		Subject:      subject.String(),
		NotBefore:    notBefore,
		NotAfter:     notAfter,
		SerialNumber: serialNumber.String(),
		IsCA:         dto.IsCA,
		DNSNames:     dto.DNSNames,
	}

	if parent != nil {
		metadata.Issuer = parent.Subject.String()
	}

	return certPEM, metadata, nil
}

func verifyCertificate(cert *x509.Certificate, caCert *x509.Certificate) map[string]interface{} {
	result := map[string]interface{}{
		"valid":         false,
		"expired":       false,
		"not_yet_valid": false,
		"self_signed":   false,
		"errors":        []string{},
	}

	now := time.Now()

	// Check expiration
	if now.Before(cert.NotBefore) {
		result["not_yet_valid"] = true
		result["errors"] = append(result["errors"].([]string), "certificate is not yet valid")
	}
	if now.After(cert.NotAfter) {
		result["expired"] = true
		result["errors"] = append(result["errors"].([]string), "certificate has expired")
	}

	// Check if self-signed
	if cert.Issuer.String() == cert.Subject.String() {
		result["self_signed"] = true
		// Verify self-signed signature
		err := cert.CheckSignatureFrom(cert)
		if err != nil {
			result["errors"] = append(result["errors"].([]string), fmt.Sprintf("invalid self-signature: %s", err.Error()))
		}
	}

	// If CA cert provided, verify against it
	if caCert != nil {
		err := cert.CheckSignatureFrom(caCert)
		if err != nil {
			result["errors"] = append(result["errors"].([]string), fmt.Sprintf("signature verification failed: %s", err.Error()))
		} else {
			result["valid"] = true
		}
	} else if result["self_signed"].(bool) && len(result["errors"].([]string)) == 0 {
		result["valid"] = true
	}

	// Add certificate details
	result["subject"] = cert.Subject.String()
	result["issuer"] = cert.Issuer.String()
	result["not_before"] = cert.NotBefore
	result["not_after"] = cert.NotAfter
	result["serial_number"] = cert.SerialNumber.String()
	result["is_ca"] = cert.IsCA
	result["dns_names"] = cert.DNSNames

	return result
}
