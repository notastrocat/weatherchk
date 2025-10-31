package main

import (
	"fmt"
	"context"
	"time"
    "encoding/json"

	"github.com/redis/go-redis/v9"
    "github.com/nitishm/go-rejson/v4"
)

var ctx = context.Background()

// Creates and returns a new Redis client
func NewRedisClient() (*redis.Client) {
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", // Redis server address
        Password: "",               // No password set
        DB:       0,                // Default DB
    })

	return rdb
}

// ReJSONClient initializes a ReJSON handler with the given Redis client and sets it up.
func ReJSONClient(rdb *redis.Client) *rejson.Handler {
    rh := rejson.NewReJSONHandler()
	rh.SetGoRedisClientWithContext(ctx, rdb)

    return rh
}

// soon to be deprecated
func SetData(rdb *redis.Client, key, val string) bool {
    // Set a value in Redis
    err := rdb.SetEx(ctx, key, val, 6 * time.Hour).Err()
    if err != nil {
        // fmt.Println(ErrStyle.Render(fmt.Sprintf("❌ Couldn't set the key: %v", err)))
		return false
    }
    // fmt.Println(SuccessStyle.Render("✔ Key '%v' set successfully!", key))
	return true
}

// soon to be deprecated
func GetVal(rdb *redis.Client, key string) {
    // Get the value from Redis
    val, err := rdb.Get(ctx, key).Result()
    if err != nil {
        fmt.Println(ErrStyle.Render(fmt.Sprintf("❌ Couldn't get the key: %v", err)))
        return
    }
    fmt.Println(SuccessStyle.Render(fmt.Sprintf("✔ Key '%v' has value: %v", key, val)))
}

// Sets the weather data in redis using rejson. The key is the city name.
func SetWeatherData(rh *rejson.Handler, key string, val map[string]interface{}) (bool, error) {
	res, err := rh.JSONSet(key, ".", val)
	if err != nil {
		return false, fmt.Errorf("rejson - Failed to JSONSet. %v", err)
	}

	if res.(string) == "OK" {
		fmt.Printf("Success: %s\n", res)
	} else {
		fmt.Printf("rejson - Failed to Set: %s\n", res)
	}

    return true, nil
}

// Gets the weather data from redis using rejson. The key is the city name.
func GetWeatherData(rh *rejson.Handler, key string) {
	res, err := rh.JSONGet(key, ".")
	if err != nil {
		fmt.Println(ErrStyle.Render(fmt.Sprintf("❌ rejson - Failed to JSONGet: %v", err)))
		return
	}

	// Convert the result to []byte
	var weatherJSON []byte
	switch v := res.(type) {
	case []byte:
		weatherJSON = v
	case string:
		weatherJSON = []byte(v)
	default:
		fmt.Println(ErrStyle.Render(fmt.Sprintf("❌ rejson - Unexpected type from JSONGet: %T", res)))
		return
	}

	readWeather := make(map[string]interface{})
	err = json.Unmarshal(weatherJSON, &readWeather)
	if err != nil {
		fmt.Println(ErrStyle.Render(fmt.Sprintf("❌ failed to JSON Unmarshal: %v", err)))
		return
	}

	fmt.Printf("Weather read from redis : %#v\n", (readWeather["currentConditions"].(map[string]interface{})["temp"]))
}