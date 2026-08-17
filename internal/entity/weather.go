package entity

import "math"

type Weather struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func NewWeather(tempC float64) Weather {
	return Weather{
		TempC: round2(tempC),
		TempF: round2(tempC*1.8 + 32),
		TempK: round2(tempC + 273.15),
	}
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}
