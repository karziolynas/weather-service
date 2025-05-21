package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type PageData struct {
	Title       string
	City        string
	Temperature string
	FeelsLike   string
	WindSpeed   string
}

func handleWeather(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getWeather(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func handleCityWeather(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		//was used for testing
		//http.Error(w, "Test", http.StatusForbidden)
		getCityWeather(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

	}
}

func getWeather(w http.ResponseWriter, r *http.Request) {
	var cityName = r.URL.Query().Get("city")
	var countryCode = r.URL.Query().Get("countryCode")
	if cityName == "" {
		http.Error(w, "No city specified", http.StatusNoContent)
	}
	if countryCode == "" {
		http.Error(w, "No country code specified", http.StatusNoContent)
	}

	data, err := GetLongLat(cityName, countryCode, http.Client{})
	if err != nil {
		http.Error(w, "Error getting area data", http.StatusNoContent)
	}
	weatherData, err := GetWeatherData(data.Lat, data.Lon, http.Client{})
	if err != nil {
		http.Error(w, "Error getting area weather", http.StatusNoContent)
	}

	defer r.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(weatherData)
}

func getCityWeather(w http.ResponseWriter, r *http.Request) {
	data, err := GetLongLat(os.Args[1], "LT", http.Client{})
	if err != nil {
		http.Error(w, "Error getting area data", http.StatusInternalServerError)
	}
	weatherData, err := GetWeatherData(data.Lat, data.Lon, http.Client{})
	if err != nil {
		http.Error(w, "Error getting area weather", http.StatusInternalServerError)
	}

	pageData := PageData{
		Title:       "City weather",
		City:        os.Args[1],
		Temperature: fmt.Sprintf("%.1f", weatherData.Main.Temp),
		FeelsLike:   fmt.Sprintf("%.1f", weatherData.Main.FeelsLike),
		WindSpeed:   fmt.Sprintf("%.1f", weatherData.Wind.Speed),
	}

	defer r.Body.Close()

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	tmpl.Execute(w, pageData)
}
