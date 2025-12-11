package secretssdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Certificate struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"created_at"`
	Name      string  `json:"name"`
	CertType  string  `json:"cert_type"`
	Algorithm string  `json:"algorithm"`
	KeySize   *int64  `json:"key_size"`
	PemData   string  `json:"pem_data"`
	Metadata  *string `json:"metadata"`
}

type GenerateKeyPairRequest struct {
	Name      string `json:"name"`
	Algorithm string `json:"algorithm"` // "RSA", "ECDSA", "ED25519"
	KeySize   *int   `json:"key_size"`
}

type GenerateKeyPairResponse struct {
	PrivateKey Certificate `json:"private_key"`
	PublicKey  Certificate `json:"public_key"`
}

type ImportCertificateRequest struct {
	Name     string `json:"name"`
	CertType string `json:"cert_type"`
	PemData  string `json:"pem_data"`
}

type GenerateCertificateRequest struct {
	Name            string   `json:"name"`
	PrivateKeyName  string   `json:"private_key_name"`
	Subject         Subject  `json:"subject"`
	ValidityDays    int      `json:"validity_days"`
	IsCA            bool     `json:"is_ca"`
	DNSNames        []string `json:"dns_names,omitempty"`
	EmailAddresses  []string `json:"email_addresses,omitempty"`
	SigningCertName *string  `json:"signing_cert_name,omitempty"`
}

type Subject struct {
	CommonName         string `json:"common_name"`
	Organization       string `json:"organization,omitempty"`
	OrganizationalUnit string `json:"organizational_unit,omitempty"`
	Country            string `json:"country,omitempty"`
	Province           string `json:"province,omitempty"`
	Locality           string `json:"locality,omitempty"`
}

type VerifyCertificateRequest struct {
	CertificateName string  `json:"certificate_name"`
	CACertName      *string `json:"ca_cert_name,omitempty"`
}

type ExportCertificateResponse struct {
	Name    string `json:"name"`
	PemData string `json:"pem_data"`
}

// ListCertificates retrieves all certificates
func (c *Client) ListCertificates() ([]Certificate, error) {
	endpoint := fmt.Sprintf("%s/api/certificates", c.BaseUrl)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.GetHttpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Success bool          `json:"success"`
		Data    []Certificate `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to list certificates")
	}

	return result.Data, nil
}

// GetCertificate retrieves a specific certificate by name
func (c *Client) GetCertificate(name string) (*Certificate, error) {
	endpoint := fmt.Sprintf("%s/api/certificates/%s", c.BaseUrl, name)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.GetHttpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Success bool        `json:"success"`
		Data    Certificate `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to get certificate")
	}

	return &result.Data, nil
}

// GenerateKeyPair generates a new key pair (private and public key)
func (c *Client) GenerateKeyPair(req GenerateKeyPairRequest) (*GenerateKeyPairResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/certificates/generate-keypair", c.BaseUrl)
	httpReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.GetHttpClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Success bool                    `json:"success"`
		Data    GenerateKeyPairResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to generate key pair")
	}

	return &result.Data, nil
}

// ImportCertificate imports a certificate from PEM data
func (c *Client) ImportCertificate(req ImportCertificateRequest) (*Certificate, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/certificates/import", c.BaseUrl)
	httpReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.GetHttpClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Success bool        `json:"success"`
		Data    Certificate `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to import certificate")
	}

	return &result.Data, nil
}

// ExportCertificate exports a certificate as PEM data
func (c *Client) ExportCertificate(name string) (*ExportCertificateResponse, error) {
	endpoint := fmt.Sprintf("%s/api/certificates/%s/export", c.BaseUrl, name)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.GetHttpClient().Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Success bool                      `json:"success"`
		Data    ExportCertificateResponse `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to export certificate")
	}

	return &result.Data, nil
}

// GenerateCertificate generates a new X.509 certificate
func (c *Client) GenerateCertificate(req GenerateCertificateRequest) (*Certificate, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/certificates/generate-certificate", c.BaseUrl)
	httpReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.GetHttpClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Success bool        `json:"success"`
		Data    Certificate `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to generate certificate")
	}

	return &result.Data, nil
}

// VerifyCertificate verifies a certificate's validity and signature
func (c *Client) VerifyCertificate(req VerifyCertificateRequest) (map[string]interface{}, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/certificates/verify", c.BaseUrl)
	httpReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.GetHttpClient().Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("failed to verify certificate")
	}

	return result.Data, nil
}

// DeleteCertificate deletes a certificate by name
func (c *Client) DeleteCertificate(name string) error {
	endpoint := fmt.Sprintf("%s/api/certificates/%s", c.BaseUrl, name)
	req, err := http.NewRequest(http.MethodDelete, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.GetHttpClient().Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	var result struct {
		Success bool `json:"success"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("failed to delete certificate")
	}

	return nil
}
