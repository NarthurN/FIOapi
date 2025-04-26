package gender

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

type Gender struct {
	Gender string `json:"gender"`
}

func (c *Client) GetGender(ctx context.Context, name string) (string, error) {
	op := "internal/apiclients/gender/gender.go.GetGender"
	params := url.Values{}
	params.Add("name", name)

	urlInstance, err := url.Parse(c.baseURL)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	urlInstance.RawQuery = params.Encode()
	formattedURL := urlInstance.String()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, formattedURL, nil)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("неожиданный статус ответа: %d", resp.StatusCode)
	}

	var prediction Gender
	if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return prediction.Gender, nil
}
