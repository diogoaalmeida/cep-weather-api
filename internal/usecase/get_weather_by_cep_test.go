package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/diogoaalmeida/cep-weather-api/internal/entity"
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

func TestGivenValidCEP_WhenExecute_ThenShouldReturnConvertedTemperatures(t *testing.T) {
	uc := usecase.NewGetWeatherByCEPUseCase(
		&locationFinderStub{location: usecase.Location{City: "Sao Paulo"}},
		&weatherFinderStub{tempC: 28.5},
	)

	weather, err := uc.Execute(context.Background(), "01310100")

	assert.NoError(t, err)
	assert.Equal(t, entity.Weather{TempC: 28.5, TempF: 83.3, TempK: 301.65}, weather)
}

func TestGivenMalformedCEP_WhenExecute_ThenShouldReturnInvalidZipcodeWithoutCallingFinders(t *testing.T) {
	uc := usecase.NewGetWeatherByCEPUseCase(
		&locationFinderStub{err: errors.New("should not be called")},
		&weatherFinderStub{err: errors.New("should not be called")},
	)

	_, err := uc.Execute(context.Background(), "123")

	assert.ErrorIs(t, err, entity.ErrInvalidZipcode)
}

func TestGivenUnknownCEP_WhenExecute_ThenShouldReturnZipcodeNotFound(t *testing.T) {
	uc := usecase.NewGetWeatherByCEPUseCase(
		&locationFinderStub{err: entity.ErrZipcodeNotFound},
		&weatherFinderStub{tempC: 28.5},
	)

	_, err := uc.Execute(context.Background(), "00000000")

	assert.ErrorIs(t, err, entity.ErrZipcodeNotFound)
}

func TestGivenWeatherLookupFails_WhenExecute_ThenShouldPropagateError(t *testing.T) {
	weatherErr := errors.New("weather api unavailable")
	uc := usecase.NewGetWeatherByCEPUseCase(
		&locationFinderStub{location: usecase.Location{City: "Sao Paulo"}},
		&weatherFinderStub{err: weatherErr},
	)

	_, err := uc.Execute(context.Background(), "01310100")

	assert.ErrorIs(t, err, weatherErr)
}
