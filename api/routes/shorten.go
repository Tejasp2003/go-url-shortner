package routes

import (
	"os"
	"strconv"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/tejasp2003/go-url-shortner/database"
	"github.com/tejasp2003/go-url-shortner/helpers"
	"github.com/asaskevich/govalidator"
)

// request represents the expected JSON payload for shortening a URL.
type request struct {
	URL         string        `json:"url"`
	CustomShort string        `json:"short"`
	Expiry      time.Duration `json:"expiry"`
}

// response represents the JSON payload returned after a URL is shortened.
type response struct {
	URL             string        `json:"url"`
	CustomShort     string        `json:"custom_short"`
	Expiry          time.Duration `json:"expiry"`
	XRateRemaining  int           `json:"rate_limit"`
	XRateLimitReset time.Duration `json:"rate_limit_reset"`
}

// ShortenURL handles the shortening of a URL.
func ShortenURL(c *fiber.Ctx) error {
	// Parse the request body into the `request` struct
	body := new(request)
	if err := c.BodyParser(body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot parse JSON",
		})
	}

	// Implement rate limiting
	rateLimitClient := database.CreateClient(1) // Redis client for rate limiting
	defer rateLimitClient.Close()

	ipAddress := c.IP()
	rateLimitValue, err := rateLimitClient.Get(database.Ctx, ipAddress).Result()

	if err == redis.Nil {
		// Set the initial rate limit value if it's the first request from this IP
		rateLimitClient.Set(database.Ctx, ipAddress, os.Getenv("API_QUOTA"), 30*time.Minute).Err()
	} else {
		// Check remaining quota
		remainingQuota, _ := strconv.Atoi(rateLimitValue)
		if remainingQuota <= 0 {
			ttl, _ := rateLimitClient.TTL(database.Ctx, ipAddress).Result() // Get the time until the rate limit resets
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"error":            "Rate Limit Exceeded",
				"rate_limit_reset": ttl / time.Minute,
			})
		}
	}

	// Check if the input is a valid URL
	if !govalidator.IsURL(body.URL) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid URL",
		})
	}

	// Check for domain error
	if !helpers.RemoveDomainError(body.URL) {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "invalid URL",
		})
	}

	// Enforce HTTP/HTTPS
	body.URL = helpers.EnforceHTTP(body.URL)

	// Determine the short ID for the URL
	var shortID string
	if body.CustomShort == "" {
		// Generate a new short ID if not provided
		shortID = uuid.New().String()[0:6]
	} else {
		shortID = body.CustomShort
	}

	urlClient := database.CreateClient(0) // Redis client for storing URLs
	defer urlClient.Close()

	// Check if the short ID already exists
	existingURL, _ := urlClient.Get(database.Ctx, shortID).Result()
	if existingURL != "" {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": "Short URL already exists",
		})
	}

	// Set the expiry time for the short URL
	if body.Expiry == 0 {
		body.Expiry = 24 * time.Hour
	} else {
		body.Expiry = body.Expiry * time.Hour
	}

	// Store the short URL in Redis with the specified expiry time
	err = urlClient.Set(database.Ctx, shortID, body.URL, body.Expiry).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "unable to connect to the server",
		})
	}

	// Decrease the remaining quota for the user's IP
	rateLimitClient.Decr(database.Ctx, ipAddress)
	rateLimitValue, _ = rateLimitClient.Get(database.Ctx, ipAddress).Result()
	remainingQuota, _ := strconv.Atoi(rateLimitValue)
	ttl, _ := rateLimitClient.TTL(database.Ctx, ipAddress).Result()

	// Prepare the response
	resp := response{
		URL:             body.URL,
		CustomShort:     os.Getenv("DOMAIN") + "/" + shortID,
		Expiry:          body.Expiry,
		XRateRemaining:  remainingQuota,
		XRateLimitReset: ttl / time.Minute,
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}
