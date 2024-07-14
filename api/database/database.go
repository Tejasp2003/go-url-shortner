package database

import (
	"context"
	"github.com/redis/go-redis/v9"
	"os"
)

// Create a global context that can be used across the application.
// Context is used for managing deadlines, cancelation signals, and other request-scoped values.
var Ctx = context.Background()

// CreateClient initializes and returns a Redis client.
// dbNo: The Redis database number to connect to (e.g., 0 for the default database).
func CreateClient(dbNo int) *redis.Client {
	// Create a new Redis client with the specified options.
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("DB_ADDR"), // Redis server address, retrieved from environment variables.
		Password: os.Getenv("DB_PASS"), // Redis server password, retrieved from environment variables.
		DB:       dbNo,                 // The database number to connect to.
	})

	// Return the created Redis client.
	return rdb
}
