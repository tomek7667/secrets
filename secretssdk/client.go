package secretssdk

type Client struct {
	BaseUrl string
	Token   string
}

func New(baseUrl, token string) (*Client, error) {
	c := &Client{
		BaseUrl: baseUrl,
		Token:   token,
	}
	if err := c.Ping(); err != nil {
		return nil, err
	}
	return c, nil
}
