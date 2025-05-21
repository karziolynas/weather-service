package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
)

var tmpl *template.Template

func SetupHTTP(portNumber int) {
	tmpl = template.Must(template.ParseFiles("/go/src/weather_service/index.html"))

	http.HandleFunc("/weather", handleWeather)

	http.HandleFunc("/cityWeather", handleCityWeather)

	log.Printf("Server started on port %v \n", portNumber+1)
	address := "0.0.0.0:" + strconv.Itoa(portNumber+1)
	log.Fatal(http.ListenAndServe(address, nil))

}
