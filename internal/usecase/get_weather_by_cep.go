package usecase

import (
	"context"

	"github.com/diogoaalmeida/cep-weather-api/internal/entity"
)

type Location struct {
	City string
}

type LocationFinder interface {
	Find(ctx context.Context, cep string) (Location, error)
}

type WeatherFinder interface {
	FindByCity(ctx context.Context, city string) (tempC float64, err error)
}

type GetWeatherByCEPUseCase struct {
	LocationFinder LocationFinder
	WeatherFinder  WeatherFinder
}

func NewGetWeatherByCEPUseCase(locationFinder LocationFinder, weatherFinder WeatherFinder) *GetWeatherByCEPUseCase {
	return &GetWeatherByCEPUseCase{
		LocationFinder: locationFinder,
		WeatherFinder:  weatherFinder,
	}
}

func (uc *GetWeatherByCEPUseCase) Execute(ctx context.Context, cep string) (entity.Weather, error) {
	if !entity.IsValidCEP(cep) {
		return entity.Weather{}, entity.ErrInvalidZipcode
	}

	location, err := uc.LocationFinder.Find(ctx, cep)
	if err != nil {
		return entity.Weather{}, err
	}

	tempC, err := uc.WeatherFinder.FindByCity(ctx, location.City)
	if err != nil {
		return entity.Weather{}, err
	}

	return entity.NewWeather(tempC), nil
}
