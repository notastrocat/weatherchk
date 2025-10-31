package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"
)

type Client struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

var unitGroup = []string{"metric", "us", "uk"}

func WeatherClient(apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{},
		apiKey:     apiKey,
		baseURL:    "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/",
	}
}

func (c *Client) GetCurrentWeather() (map[string]interface{}, time.Duration, error) {
	city := LocationInput
	unitType := UnitTypeInput
	if city == "" {
		return nil, 0, fmt.Errorf("city cannot be empty")
	}
	if unitType == "" {
		fmt.Println(WarnStyle.Render("⚠ No unit type provided, defaulting to 'metric'"))
		unitType = "metric"
	} else {
		if !( slices.Contains(unitGroup, unitType) ) {
			fmt.Println(WarnStyle.Render("⚠ invalid unit type provided. Valid options are: metric/us/uk. Defaulting to 'metric'"))
		}
		unitType = "metric"
	}

	url := fmt.Sprintf("%s%s?unitGroup=%s&key=%s&contentType=json", c.baseURL, city, unitType, c.apiKey)
	// fmt.Printf(WarnStyle.Render("🤔 Fetching weather data from %s\n\n"), url)
	// fmt.Printf(WarnStyle.Render("🤔 api key: %s\n"), c.apiKey)
	// fmt.Printf(WarnStyle.Render("🤔 Fetching weather data for %s with unit type %s...\n"), city, unitType)

	// url = "https://weather.visualcrossing.com/VisualCrossingWebServices/rest/services/timeline/Jaipur?unitGroup=uk&key=MVJ2XRJEGFPXG3FAP6SCPTCGZ&contentType=json"

	var startTime = time.Now()
	resp, err := c.httpClient.Get(url)
	var timeTaken = time.Since(startTime)
	if err != nil {
		return nil, timeTaken, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, timeTaken, fmt.Errorf("failed to get weather data: %s", resp.Status)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, timeTaken, err
	}

	// Extract current conditions (this is a simplified example)
	// currentConditions, ok := result["currentConditions"].(map[string]interface{})
	// if !ok {
	// 	var timeTaken = time.Since(startTime)
	// 	return "", timeTaken, fmt.Errorf("invalid response format")
	// }

	// temperature := currentConditions["temp"].(float64)
	// conditions := currentConditions["conditions"].(string)

	return result, timeTaken, nil
}
