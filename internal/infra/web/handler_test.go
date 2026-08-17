package web_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diogoaalmeida/cep-weather-api/internal/entity"
	"github.com/diogoaalmeida/cep-weather-api/internal/infra/web"
	"github.com/diogoaalmeida/cep-weather-api/internal/usecase"
	"github.com/stretchr/testify/assert"
)

type locationFinderStub struct {
	location usecase.Location
	err      error
}

func (s *locationFinderStub) Find(_ context.Context, _ string) (usecase.Location, error) {
	return s.location, s.err
}

type weatherFinderStub struct {
	tempC float64
	err   error
}

func (s *weatherFinderStub) FindByCity(_ context.Context, _ string) (float64, error) {
	return s.tempC, s.err
}

func newTestHandler(locationErr, weatherErr error, tempC float64) *web.WeatherHandler {
	uc := usecase.NewGetWeatherByCEPUseCase(
		&locationFinderStub{location: usecase.Location{City: "Sao Paulo"}, err: locationErr},
		&weatherFinderStub{tempC: tempC, err: weatherErr},
	)
	return web.NewWeatherHandler(uc)
}

func doRequest(t *testing.T, handler *web.WeatherHandler, cep string) *httptest.ResponseRecorder {
	t.Helper()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/{cep}", handler.Get)

	req := httptest.NewRequest(http.MethodGet, "/weather/"+cep, nil)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, req)
	return rec
}

func TestGivenValidCEP_WhenGet_ThenShouldReturn200WithTemperatures(t *testing.T) {
	handler := newTestHandler(nil, nil, 28.5)

	rec := doRequest(t, handler, "01310100")

	assert.Equal(t, http.StatusOK, rec.Code)

	var body entity.Weather
	assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	assert.Equal(t, entity.Weather{TempC: 28.5, TempF: 83.3, TempK: 301.65}, body)
}

func TestGivenMalformedCEP_WhenGet_ThenShouldReturn422(t *testing.T) {
	handler := newTestHandler(nil, nil, 28.5)

	rec := doRequest(t, handler, "123")

	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
	assert.JSONEq(t, `{"message":"invalid zipcode"}`, rec.Body.String())
}

func TestGivenUnknownCEP_WhenGet_ThenShouldReturn404(t *testing.T) {
	handler := newTestHandler(entity.ErrZipcodeNotFound, nil, 0)

	rec := doRequest(t, handler, "00000000")

	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.JSONEq(t, `{"message":"can not find zipcode"}`, rec.Body.String())
}
