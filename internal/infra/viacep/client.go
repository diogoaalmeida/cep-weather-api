package viacep

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/diogoaalmeida/cep-weather-api/internal/entity"
	"github.com/diogoaalmeida/cep-weather-api/internal/usecase"
)

const defaultBaseURL = "https://viacep.com.br/ws"

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewClient() *Client {
	return &Client{
		BaseURL:    defaultBaseURL,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}
}

type response struct {
	Localidade string `json:"localidade"`
}

// ViaCEP's erro field is typed inconsistently (bool or string), so we
// check for an empty city name instead of parsing erro.
func (c *Client) Find(ctx context.Context, cep string) (usecase.Location, error) {
	url := fmt.Sprintf("%s/%s/json/", c.BaseURL, cep)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return usecase.Location{}, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return usecase.Location{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return usecase.Location{}, fmt.Errorf("viacep: unexpected status %d", resp.StatusCode)
	}

	var r response
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return usecase.Location{}, err
	}

	if r.Localidade == "" {
		return usecase.Location{}, entity.ErrZipcodeNotFound
	}

	return usecase.Location{City: r.Localidade}, nil
}
