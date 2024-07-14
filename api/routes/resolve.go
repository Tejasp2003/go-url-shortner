package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/tejasp2003/go-url-shortner/database"
)

// ResolveURL handles the redirection from a shortened URL to the original URL.
func ResolveURL(c *fiber.Ctx) error {
	// Extract the shortened URL parameter from the request context.
	shortURL := c.Params("url")

	// Create a Redis client for interacting with the URL storage (database 0).
	urlClient := database.CreateClient(0)
	defer urlClient.Close()

	// Attempt to retrieve the original URL from Redis using the shortened URL.
	originalURL, err := urlClient.Get(database.Ctx, shortURL).Result()
	if err == redis.Nil {
		// If the shortened URL is not found, return a 404 status with an error message.
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Short URL not found",
		})
	} else if err != nil {
		// If there is any other error, return a 500 status with an error message.
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal Server Error",
		})
	}

	// Create a new Redis client for incrementing the access count (database 1).
	countClient := database.CreateClient(1)
	defer countClient.Close()

	// Increment the access count for the shortened URL.
	countClient.Incr(database.Ctx, shortURL+":count")

	// Redirect the user to the original URL with a 301 status (Moved Permanently).
	return c.Redirect(originalURL, fiber.StatusMovedPermanently)
}
