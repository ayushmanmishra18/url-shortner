package main

import (
	"net/url"
	"sync"

	"github.com/gofiber/fiber/v2"
)

var urlStore = make(map[string]string)

// mutex for thread-safe counter

var mu sync.Mutex
var counter int64 = 1

const base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Base62 encoder to generate short codes
func encodeBase62(num int64) string {
	if num == 0 {
		return string(base62Chars[0])
	}

	var encoded []byte

	for num > 0 {
		remainder := num % 62
		encoded = append([]byte{base62Chars[remainder]}, encoded...)
		num = num / 62
	}

	return string(encoded)
}

// validate URL
func isValidURL(input string) bool {
	u, err := url.ParseRequestURI(input)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func main() {
	app := fiber.New()

	// health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Server is running",
		})
	})

	// shorten URL
	app.Post("/shorten", func(c *fiber.Ctx) error {
		type Request struct {
			URL string `json:"url"`
		}

		var body Request
		if err := c.BodyParser(&body); err != nil || body.URL == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		// validate URL
		if !isValidURL(body.URL) {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid URL format",
			})
		}

		// thread-safe counter increment
		mu.Lock()
		code := encodeBase62(counter)
		counter++
		mu.Unlock()

		urlStore[code] = body.URL

		return c.JSON(fiber.Map{
			"short_url": "http://localhost:3000/" + code,
		})
	})

	// redirect
	app.Get("/:code", func(c *fiber.Ctx) error {
		code := c.Params("code")

		url, exists := urlStore[code]
		if !exists {
			return c.Status(404).JSON(fiber.Map{
				"error": "URL not found",
			})
		}

		return c.Redirect(url, 302)
	})

	app.Listen(":3000")
}