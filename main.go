package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/karziolynas/goconsul"
)

const serviceNameConst string = "weatherService"
const address string = "127.0.0.1"

//const port int = 3212

type WeatherData struct {
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func GetLongLat(city string, countryCode string, c http.Client) (*Location, error) {
	req, err := http.NewRequest("GET", "http://api.openweathermap.org/geo/1.0/direct", nil)
	if err != nil {
		return nil, err
	}

	fullLocation := city + "," + countryCode
	queryParams := req.URL.Query()
	queryParams.Set("q", fullLocation)
	queryParams.Set("limit", "3")
	queryParams.Set("appid", "39dcd580770065f123a2d383a1b78650")
	req.URL.RawQuery = queryParams.Encode()

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []Location
	errJson := json.NewDecoder(resp.Body).Decode(&data)
	if errJson != nil {
		return nil, errJson
	}

	return &data[0], nil
}

func GetWeatherData(lat float64, lon float64, c http.Client) (*WeatherData, error) {
	req, err := http.NewRequest("GET", "https://api.openweathermap.org/data/2.5/weather?", nil)
	if err != nil {
		return nil, err
	}

	queryParams := req.URL.Query()
	queryParams.Set("lat", strconv.FormatFloat(lat, 'f', -1, 32))
	queryParams.Set("lon", strconv.FormatFloat(lon, 'f', -1, 32))
	queryParams.Set("units", "metric")
	queryParams.Set("appid", "39dcd580770065f123a2d383a1b78650")
	req.URL.RawQuery = queryParams.Encode()

	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data WeatherData
	errJson := json.NewDecoder(resp.Body).Decode(&data)
	if errJson != nil {
		return nil, errJson
	}

	return &data, nil
}

func main() {

	client := http.Client{
		Timeout: 2 * time.Minute,
	}

	data, err := GetLongLat(os.Args[1], "LT", client)
	if err != nil {
		log.Fatal(err)
	}

	portNumber, err := strconv.Atoi(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	// makes the tags and parameters
	consulAddr := /*"127.0.0.1:8500"*/ "host.docker.internal:8501"
	tags := make([]string, 2)
	tags[0] = "weather_" + os.Args[1]
	tags[1] = "test"
	serviceName := serviceNameConst + "_" + os.Args[1]
	service := goconsul.NewService(consulAddr, os.Args[2], serviceName, address, portNumber, tags)

	portNumberService := strconv.Itoa(portNumber + 1)
	fullServiceAddress := "http://" + address + ":" + portNumberService + "/cityWeather"
	//fullServiceAddress := "http://host.docker.internal:" + portNumberService + "/cityWeather"

	go SetupHTTP(portNumber)

	go service.Start(consulAddr, fullServiceAddress)

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		weatherData, err := GetWeatherData(data.Lat, data.Lon, client)
		if err != nil {
			log.Println("Error fetching weather data:", err)
		} else {
			fmt.Println("|========WEATHER DATA===================================|")
			fmt.Printf("|  City: %s \n", os.Args[1])
			fmt.Printf("|  Temp: %.2f°C \n", weatherData.Main.Temp)
			fmt.Printf("|  Feels like: %.2f°C \n", weatherData.Main.FeelsLike)
			fmt.Printf("|  Wind speed: %.2f m/s \n", weatherData.Wind.Speed)
			fmt.Println("|=======================================================|")

		}
		<-ticker.C
	}

}
