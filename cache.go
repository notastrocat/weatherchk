package main

import (
	"fmt"
	"context"
	"time"

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


    // separate function to get the weather data
	// studentJSON, err := redis.Bytes(rh.JSONGet("student", "."))
	// if err != nil {
	// 	fmt.Errorf("rejson - Failed to JSONGet. %v", err)
	// 	return
	// }

	// readStudent := Student{}
	// err = json.Unmarshal(studentJSON, &readStudent)
	// if err != nil {
	// 	fmt.Errorf("Failed to JSON Unmarshal")
	// 	return
	// }

	// fmt.Printf("Student read from redis : %#v\n", readStudent)
