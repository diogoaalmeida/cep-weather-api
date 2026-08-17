package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/diogoaalmeida/cep-weather-api/internal/entity"
	"github.com/diogoaalmeida/cep-weather-api/internal/usecase"
)

type WeatherHandler struct {
	UseCase *usecase.GetWeatherByCEPUseCase
}

func NewWeatherHandler(useCase *usecase.GetWeatherByCEPUseCase) *WeatherHandler {
	return &WeatherHandler{UseCase: useCase}
}

func (h *WeatherHandler) Get(w http.ResponseWriter, r *http.Request) {
	cep := r.PathValue("cep")

	weather, err := h.UseCase.Execute(r.Context(), cep)
	switch {
	case errors.Is(err, entity.ErrInvalidZipcode):
		respondError(w, http.StatusUnprocessableEntity, err.Error())
	case errors.Is(err, entity.ErrZipcodeNotFound):
		respondError(w, http.StatusNotFound, err.Error())
	case err != nil:
		log.Printf("get weather by cep %q: %v", cep, err)
		respondError(w, http.StatusInternalServerError, "internal server error")
	default:
		respondJSON(w, http.StatusOK, weather)
	}
}

func respondJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, map[string]string{"message": message})
}
