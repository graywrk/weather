package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	// Shortening the import reference name seems to make it a bit easier
	owm "github.com/briandowns/openweathermap"
)

var apiKey = os.Getenv("OWM_API_KEY")

const IP_TO_COORDINATE_URL = "http://ip-api.com/json"

type DataFromIpApi struct {
	Status      string  `json:"status"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	Region      string  `json:"region"`
	RegionName  string  `json:"regionName"`
	City        string  `json:"city"`
	Zip         string  `json:"zip"`
	Lat         float64 `json:"lat"`
	Lon         float64 `json:"lon"`
	Timezone    string  `json:"timezone"`
	ISP         string  `json:"isp"`
	ORG         string  `json:"org"`
	AS          string  `json:"as"`
	Message     string  `json:"message"`
	Query       string  `json:"query"`
}

func getLocation() (*DataFromIpApi, error) {
	response, err := http.Get(IP_TO_COORDINATE_URL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	result, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	r := &DataFromIpApi{}
	if err = json.Unmarshal(result, &r); err != nil {
		return nil, err
	}
	return r, nil
}

// getCurrent gets the current weather for the provided
// location in the units provided.
func getCurrent(location *owm.Coordinates, units, lang string) (*owm.CurrentWeatherData, error) {
	w, err := owm.NewCurrent(units, lang, os.Getenv("OWM_API_KEY"))
	if err != nil {
		return nil, err
	}
	w.CurrentByCoordinates(location)
	return w, nil
}

func getForecast5(location *owm.Coordinates, units, lang string) (*owm.Forecast5WeatherData, error) {
	w, err := owm.NewForecast("5", units, lang, os.Getenv("OWM_API_KEY"))
	if err != nil {
		return nil, err
	}
	w.DailyByCoordinates(location, 5)
	forecast := w.ForecastWeatherJson.(*owm.Forecast5WeatherData)
	return forecast, err
}

func forecastString(w *owm.Forecast5WeatherList) string {
	time := fmt.Sprintf("%v", w.DtTxt)
	currentTemp := fmt.Sprintf("%.1f", w.Main.Temp)
	minTemp := fmt.Sprintf("%.1f", w.Main.TempMax)
	maxTemp := fmt.Sprintf("%.1f", w.Main.TempMin)
	feelsLikeTemp := fmt.Sprintf("%.0f", w.Main.FeelsLike)
	humidity := fmt.Sprintf("%d", w.Main.Humidity)
	wind := fmt.Sprintf("%.1f", w.Wind.Speed)
	rain := fmt.Sprintf("%.0f", w.Rain.OneH)
	snow := fmt.Sprintf("%.0f", w.Snow.OneH)
	return time + ": 🌡️ " + currentTemp + "℃ (от " + minTemp + "℃  до " + maxTemp + "℃ ) ощущается как " + feelsLikeTemp + "℃ 💦 " +
		humidity + "% " + "🌬️  " + wind + " м/с, ☔" + rain + "%" + " ❄️" + snow + "%"
}

func main() {
	api_key := os.Getenv("OWM_API_KEY")
	if len(api_key) == 0 {
		log.Fatalln("Задайте переменную окружения OWM_API_KEY!")
	}

	loc, err := getLocation()
	if err != nil {
		log.Fatalln(err)
	}

	w, err := getCurrent(&owm.Coordinates{
		Longitude: loc.Lon,
		Latitude:  loc.Lat},
		"C",
		"ru")
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println(loc.Country + ", " + loc.City + ", " + w.Name)

	currentTemp := fmt.Sprintf("%.1f", w.Main.Temp)
	minTemp := fmt.Sprintf("%.1f", w.Main.TempMax)
	maxTemp := fmt.Sprintf("%.1f", w.Main.TempMin)
	feelsLikeTemp := fmt.Sprintf("%.0f", w.Main.FeelsLike)
	humidity := fmt.Sprintf("%d", w.Main.Humidity)
	wind := fmt.Sprintf("%.1f", w.Wind.Speed)
	rain := fmt.Sprintf("%.0f", w.Rain.OneH)
	snow := fmt.Sprintf("%.0f", w.Snow.OneH)

	out := "Сейчас: 🌡️ " +
		currentTemp + "℃ (от " + minTemp + "℃  до " + maxTemp + "℃ ) ощущается как " + feelsLikeTemp + "℃ 💦 " +
		humidity + "% " + "🌬️  " + wind + " м/с, ☔ " + rain + "%" + " ❄️" + snow + "% "
	fmt.Println(out)

	fw, err := getForecast5(
		&owm.Coordinates{
			Longitude: loc.Lon,
			Latitude:  loc.Lat},
		"C",
		"ru")
	if err != nil {
		log.Fatalln(err)
	}

	for _, curw := range fw.List {
		fmt.Println(forecastString(&curw))

	}
}
