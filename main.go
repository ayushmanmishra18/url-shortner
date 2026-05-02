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
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})