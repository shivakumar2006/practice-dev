package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
)

type apiConfigData struct {
	OpenWeatherMapApiKey string `json:"OpenWeatherMapApiKey"`
}

type WeatherData struct {
	Name string `json:"name"`
	Main struct {
		Celcius  float64 `json:"temp"`
		Humidity float64 `json:"humidity"`
		Pressure float64 `json:"pressure"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

func loadApiConfig(filename string) (apiConfigData, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return apiConfigData{}, err
	}

	var c apiConfigData
	err = json.Unmarshal(bytes, &c)
	if err != nil {
		return apiConfigData{}, err
	}

	return c, nil
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request at /hello")
	w.Write([]byte("hello form go! \n"))
}

func query(city string) (WeatherData, error) {
	apiConfig, err := loadApiConfig(".apiConfig")
	if err != nil {
		return WeatherData{}, err
	}

	apiURl := "http://api.openweathermap.org/data/2.5/weather?APPID=" + apiConfig.OpenWeatherMapApiKey + "&q=" + city
	log.Printf("Fetching weather data for city: %s \n", city)

	response, err := http.Get(apiURl)
	if err != nil {
		return WeatherData{}, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return WeatherData{}, err
	}

	var d WeatherData
	if err := json.NewDecoder(response.Body).Decode(&d); err != nil {
		return WeatherData{}, err
	}

	d.Main.Celcius = d.Main.Celcius - 273.15

	fmt.Printf("Temperature in Celsius: %.2f°C\n", d.Main.Celcius)
	fmt.Printf("Humidity: %.2f%%\n", d.Main.Humidity)
	fmt.Printf("Pressure: %.2f hpa\n", d.Main.Pressure)
	fmt.Printf("Wind Speed: %.2f m/s\n", d.Wind.Speed)
	fmt.Printf("Weather Description: %s \n", d.Weather[0].Description)

	return d, nil
}

func main() {
	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"})

	log.Println("starting the server on the port : 8080...")
	http.HandleFunc("/hello", hello)

	http.HandleFunc("/weather/",
		func(w http.ResponseWriter, r *http.Request) {
			city := strings.SplitN(r.URL.Path, "/", 3)[2]
			data, err := query(city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json; charset=utf-8")
			json.NewEncoder(w).Encode(data)
		})
	http.ListenAndServe(":8080", handlers.CORS(origins, headers, methods)(http.DefaultServeMux))
}
