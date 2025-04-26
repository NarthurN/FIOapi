package age

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func New(baseURL string, client *http.Client) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: client,
	}
}

type Age struct {
	Age int `json:"age"`
}

func (c *Client) GetAge(ctx context.Context, name string) (int, error) {
	op := "internal/apiclients/age/age.go.GetAge"
	params := url.Values{}
	params.Add("name", name)

	urlInstance, err := url.Parse(c.baseURL)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	urlInstance.RawQuery = params.Encode()
	formattedURL := urlInstance.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, formattedURL, nil)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("неожиданный статус ответа: %d", resp.StatusCode)
	}

	var prediction Age
	if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return prediction.Age, nil
}
