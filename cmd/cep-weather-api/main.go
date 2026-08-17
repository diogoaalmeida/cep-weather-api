package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/diogoaalmeida/cep-weather-api/internal/infra/viacep"
	"github.com/diogoaalmeida/cep-weather-api/internal/infra/weatherapi"
	"github.com/diogoaalmeida/cep-weather-api/internal/infra/web"
	"github.com/diogoaalmeida/cep-weather-api/internal/usecase"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		log.Fatal("WEATHER_API_KEY is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	useCase := usecase.NewGetWeatherByCEPUseCase(viacep.NewClient(), weatherapi.NewClient(apiKey))
	handler := web.NewWeatherHandler(useCase)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /weather/{cep}", handler.Get)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
