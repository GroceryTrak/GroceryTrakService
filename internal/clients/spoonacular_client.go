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

type SpoonacularNutrient struct {
	Name                string  `json:"name"`
	Amount              float64 `json:"amount"`
	Unit                string  `json:"unit"`
	PercentOfDailyNeeds float64 `json:"percentOfDailyNeeds"`
}

type SpoonacularIngredientInfo struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Image     string `json:"image"`
	Nutrition struct {
		Nutrients []SpoonacularNutrient `json:"nutrients"`
	} `json:"nutrition"`
}

func NewSpoonacularClient(baseURL, apiKey string) *SpoonacularClient {
	return &SpoonacularClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  &http.Client{},
	}
}

func (c *SpoonacularClient) SearchIngredient(ctx context.Context, query string) (*SpoonacularIngredientInfo, error) {
	// First search for the ingredient
	searchURL := fmt.Sprintf("%s/food/ingredients/search?query=%s&number=1&apiKey=%s", c.baseURL, query, c.apiKey)

	req, err := http.NewRequestWithContext(ctx, "GET", searchURL, nil)
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

	var searchResult struct {
		Results []SpoonacularIngredient `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&searchResult); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if len(searchResult.Results) == 0 {
		return nil, fmt.Errorf("no results found")
	}

	ingredient := searchResult.Results[0]

	// Then get detailed information including nutrients
	infoURL := fmt.Sprintf("%s/food/ingredients/%d/information?apiKey=%s&amount=1", c.baseURL, ingredient.ID, c.apiKey)

	req, err = http.NewRequestWithContext(ctx, "GET", infoURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err = c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var info SpoonacularIngredientInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &info, nil
}
