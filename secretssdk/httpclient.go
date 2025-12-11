package secretssdk

import "net/http"

type authTransport struct {
	token string
	base  http.RoundTripper
}

func (t *authTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "Api "+t.token)
	return t.base.RoundTrip(req)
}

func (c *Client) GetHttpClient() *http.Client {
	return &http.Client{
		Transport: &authTransport{
			token: c.Token,
			base:  http.DefaultTransport,
		},
	}
}
