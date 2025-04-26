package nationality

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

type Nationality struct {
	Country []Country `json:"country"`
}

type Country struct {
	ID          string  `json:"country_id"`
	Probability float64 `json:"probability"`
}

func (c *Client) GetNationality(ctx context.Context, name string) (string, error) {
	op := "internal/apiclients/nationality/nationality.go.GetNationality"
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

	var prediction Nationality
	if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return prediction.Country[0].ID, nil
}
