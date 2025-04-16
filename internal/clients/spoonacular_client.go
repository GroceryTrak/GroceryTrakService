package clients

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SpoonacularClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

type SpoonacularIngredient struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Image string `json:"image"`
}

func NewSpoonacularClient(baseURL, apiKey string) *SpoonacularClient {
	return &SpoonacularClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

func (c *SpoonacularClient) SearchIngredient(ctx context.Context, query string) (*SpoonacularIngredient, error) {
	url := fmt.Sprintf("%s/food/ingredients/search?query=%s&number=1&apiKey=%s", c.baseURL, query, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result struct {
		Results []SpoonacularIngredient `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(result.Results) == 0 {
		return nil, fmt.Errorf("no results found")
	}

	return &result.Results[0], nil
}
