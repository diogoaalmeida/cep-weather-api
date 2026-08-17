# cep-weather-api

Given a Brazilian CEP, looks up the city and returns the current temperature in Celsius, Fahrenheit, and Kelvin.

## Live URL

https://cep-weather-api-893033911139.us-central1.run.app/weather/01310100

## Run locally

Copy `.env.example` to `.env` and set a [WeatherAPI](https://www.weatherapi.com/) key:

```bash
cp .env.example .env
# edit .env, set WEATHER_API_KEY
```

Run directly:

```bash
go run ./cmd/cep-weather-api
```

or via Docker:

```bash
docker build -t cep-weather-api .
docker run -p 8080:8080 --env-file .env cep-weather-api
```

## Try it

```bash
GET http://localhost:8080/weather/01310100
```

`api.http` has ready requests for a valid CEP, an invalid one, and a well-formed but nonexistent one.

## Run tests

```bash
go test ./...
```

## API

| Scenario | Status | Body |
|---|---|---|
| Valid CEP | 200 | `{"temp_C": 28.5, "temp_F": 83.3, "temp_K": 301.65}` |
| CEP is not 8 digits | 422 | `{"message": "invalid zipcode"}` |
| CEP is well-formed but not found | 404 | `{"message": "can not find zipcode"}` |

## Architecture

- `internal/entity`: `Weather` (Celsius/Fahrenheit/Kelvin conversion), CEP format validation
- `internal/usecase`: `GetWeatherByCEPUseCase`, depends on two ports (`LocationFinder`, `WeatherFinder`) so it's testable without real API calls
- `internal/infra/viacep`: resolves a CEP to a city via ViaCEP
- `internal/infra/weatherapi`: resolves a city to a current temperature via WeatherAPI
- `internal/infra/web`: maps use case errors to the status codes above
- `cmd/cep-weather-api`: wiring and startup, reads `PORT` (Cloud Run sets this) and `WEATHER_API_KEY`
