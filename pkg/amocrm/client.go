package amocrm

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type AMOClient struct {
	httpClient *http.Client
}

func (c *AMOClient) makeURL(domain, path string) string {
	return fmt.Sprintf("https://%s/%s", domain, path)
}

func NewAMOClient(httpClient *http.Client) *AMOClient {
	return &AMOClient{
		httpClient: httpClient,
	}
}

func (c *AMOClient) DoRequest(method, path, domain, accessToken string, body io.Reader) (*http.Response, error) {
	url := c.makeURL(domain, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	log.Println(url)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "amoCRM-oAuth-client/1.0")
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
