package main

import (
	"fmt"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func Connect() (*redis.Client) {
    // Create a new Redis client
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379", // Redis server address
        Password: "",               // No password set
        DB:       0,                // Default DB
    })

	return rdb
}

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

func GetVal(rdb *redis.Client, key string) {
    // Get the value from Redis
    val, err := rdb.Get(ctx, key).Result()
    if err != nil {
        fmt.Println(ErrStyle.Render(fmt.Sprintf("❌ Couldn't get the key: %v", err)))
    }
    fmt.Println(SuccessStyle.Render(fmt.Sprintf("✔ Key '%v' has value: %v", key, val)))
}
