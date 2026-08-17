package viacep_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diogoaalmeida/cep-weather-api/internal/entity"
	"github.com/diogoaalmeida/cep-weather-api/internal/infra/viacep"
	"github.com/diogoaalmeida/cep-weather-api/internal/usecase"
	"github.com/stretchr/testify/assert"
)

func newTestClient(server *httptest.Server) *viacep.Client {
	client := viacep.NewClient()
	client.BaseURL = server.URL
	return client
}

func TestGivenExistingCEP_WhenFind_ThenShouldReturnCity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"cep":"01310-100","localidade":"São Paulo","erro":false}`))
	}))
	defer server.Close()

	location, err := newTestClient(server).Find(context.Background(), "01310100")

	assert.NoError(t, err)
	assert.Equal(t, usecase.Location{City: "São Paulo"}, location)
}

func TestGivenUnknownCEP_WhenFind_ThenShouldReturnZipcodeNotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"erro":true}`))
	}))
	defer server.Close()

	_, err := newTestClient(server).Find(context.Background(), "00000000")

	assert.ErrorIs(t, err, entity.ErrZipcodeNotFound)
}
