package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type WeatherData struct {
	Location struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"location"`
	Current struct {
		TempC     float64 `json:"temp_c"`
		TempF     float64 `json:"temp_f"`
		Condition struct {
			Text string `json:"text"`
		} `json:"condition"`
		Humidity    int     `json:"humidity"`
		WindKph     float64 `json:"wind_kph"`
		FeelsLikeC  float64 `json:"feelslike_c"`
		FeelsLikeF  float64 `json:"feelslike_f"`
		UV         float64 `json:"uv"`
	} `json:"current"`
}

func getWeather(apiKey, location string) (*WeatherData, error) {
	url := fmt.Sprintf("http://api.weatherapi.com/v1/current.json?key=%s&q=%s&aqi=no", apiKey, location)
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var weather WeatherData
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, fmt.Errorf("error parsing weather data: %v", err)
	}

	return &weather, nil
}

func displayWeather(weather *WeatherData) {
	fmt.Printf("\nWeather for %s, %s\n", weather.Location.Name, weather.Location.Country)
	fmt.Println("----------------------------------------")
	fmt.Printf("Temperature: %.1f°C (%.1f°F)\n", weather.Current.TempC, weather.Current.TempF)
	fmt.Printf("Feels like: %.1f°C (%.1f°F)\n", weather.Current.FeelsLikeC, weather.Current.FeelsLikeF)
	fmt.Printf("Condition: %s\n", weather.Current.Condition.Text)
	fmt.Printf("Humidity: %d%%\n", weather.Current.Humidity)
	fmt.Printf("Wind Speed: %.1f km/h\n", weather.Current.WindKph)
	fmt.Printf("UV Index: %.1f\n", weather.Current.UV)
}

func main() {
	apiKey := flag.String("key", "", "WeatherAPI API key")
	location := flag.String("location", "", "Location to get weather for")
	flag.Parse()

	if *apiKey == "" {
		// Try to get API key from environment variable
		*apiKey = os.Getenv("WEATHER_API_KEY")
		if *apiKey == "" {
			fmt.Println("Error: API key is required. Provide it using -key flag or WEATHER_API_KEY environment variable")
			os.Exit(1)
		}
	}

	if *location == "" {
		fmt.Println("Error: Location is required. Provide it using -location flag")
		os.Exit(1)
	}

	weather, err := getWeather(*apiKey, *location)
	if err != nil {
		fmt.Printf("Error getting weather: %v\n", err)
		os.Exit(1)
	}

	displayWeather(weather)
}
