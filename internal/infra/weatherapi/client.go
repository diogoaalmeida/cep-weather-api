package weatherapi

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"
)

const defaultBaseURL = "https://api.weatherapi.com/v1"

type Client struct {
	BaseURL    string
	APIKey     string
	HTTPClient *http.Client
}

func NewClient(apiKey string) *Client {
	return &Client{
		BaseURL:    defaultBaseURL,
		APIKey:     apiKey,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

type response struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

func (c *Client) FindByCity(ctx context.Context, city string) (float64, error) {
	endpoint := fmt.Sprintf("%s/current.json?key=%s&q=%s&aqi=no",
		c.BaseURL, url.QueryEscape(c.APIKey), url.QueryEscape(city))

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return 0, err
	}

	start := time.Now()
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	log.Printf("weatherapi: city=%q status=%d duration=%s", city, resp.StatusCode, time.Since(start))

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("weatherapi: unexpected status %d", resp.StatusCode)
	}

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, err
	}

	log.Printf("weatherapi: city=%q temp_c=%.1f", city, r.Current.TempC)
	return r.Current.TempC, nil
}
