package main

import (
	"math/rand"
	"time"
	"github.com/gofiber/fiber/v2"

)
var urlStore=make(map[string]string)

//generate random short code 

func generateCode(n int )string{
	const letters="abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, n)
	for i := 0; i < n; i++ {
		code[i] = letters[rand.Intn(len(letters))]
	}
	return string(code)
}

func main() {
	app := fiber.New()

	//health check endpoint
	app.Get("/health",func(c *fiber.Ctx)error {
		return c.JSON(fiber.Map{
			"message":"Server is runnning",
		})
	})

	//shorten url

	app.Post("/shorten",func(c *fiber.Ctx)error {
		type Request struct {
			URL string `json:"url"`
		}

		
		var body Request
		if err := c.BodyParser(&body); err != nil || body.URL == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request",
			})
		}

		code := generateCode(6)
		urlStore[code] = body.URL

		return c.JSON(fiber.Map{
			"short_url": "http://localhost:3000/" + code,
		})
	})

 //redirect 
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