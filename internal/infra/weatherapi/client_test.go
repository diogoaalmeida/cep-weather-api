package weatherapi_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diogoaalmeida/cep-weather-api/internal/infra/weatherapi"
	"github.com/stretchr/testify/assert"
)

func newTestClient(server *httptest.Server) *weatherapi.Client {
	client := weatherapi.NewClient("test-key")
	client.BaseURL = server.URL
	return client
}

func TestGivenCity_WhenFindByCity_ThenShouldReturnCurrentTempC(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "test-key", r.URL.Query().Get("key"))
		assert.Equal(t, "São Paulo", r.URL.Query().Get("q"))
		w.Write([]byte(`{"current":{"temp_c":28.5}}`))
	}))
	defer server.Close()

	tempC, err := newTestClient(server).FindByCity(context.Background(), "São Paulo")

	assert.NoError(t, err)
	assert.Equal(t, 28.5, tempC)
}

func TestGivenUpstreamError_WhenFindByCity_ThenShouldReturnError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	_, err := newTestClient(server).FindByCity(context.Background(), "São Paulo")

	assert.Error(t, err)
}
