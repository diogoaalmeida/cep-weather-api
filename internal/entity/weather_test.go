package entity_test

import (
	"testing"

	"github.com/diogoaalmeida/cep-weather-api/internal/entity"
	"github.com/stretchr/testify/assert"
)

func TestGivenCelsiusTemperature_WhenNewWeather_ThenShouldConvertToFahrenheitAndKelvin(t *testing.T) {
	weather := entity.NewWeather(28.5)

	assert.Equal(t, entity.Weather{TempC: 28.5, TempF: 83.3, TempK: 301.65}, weather)
}

func TestGivenNegativeCelsiusTemperature_WhenNewWeather_ThenShouldConvertCorrectly(t *testing.T) {
	weather := entity.NewWeather(-10)

	assert.Equal(t, entity.Weather{TempC: -10, TempF: 14, TempK: 263.15}, weather)
}
