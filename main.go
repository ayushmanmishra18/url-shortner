package main

import (
	"database/sql"
	"log"
	"net/url"

	_ "github.com/lib/pq"
	"github.com/gofiber/fiber/v2"
)

var db *sql.DB

const base62Chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Base62 encoder
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

	// Supabase connection (POOLER)
	connStr := "postgres://postgres.iemvhknbvsvkhxwzhpmq:Ayushmanmishra@aws-1-ap-northeast-2.pooler.supabase.com:5432/postgres?sslmode=require"

	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Supabase DB")

	// health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Server is running",
		})
	})

	// shorten URL (DB-driven ID)
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

		if !isValidURL(body.URL) {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid URL format",
			})
		}

		var id int64

		// Step 1: Insert URL → get ID
		err := db.QueryRow(
			"INSERT INTO urls (long_url) VALUES ($1) RETURNING id",
			body.URL,
		).Scan(&id)

		if err != nil {
			log.Println("INSERT ERROR:", err)
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to insert URL",
			})
		}

		// Step 2: Generate short code from ID
		code := encodeBase62(id)

		// Step 3: Update row with short_code
		_, err = db.Exec(
			"UPDATE urls SET short_code=$1 WHERE id=$2",
			code, id,
		)

		if err != nil {
			log.Println("UPDATE ERROR:", err)
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to update short code",
			})
		}

		return c.JSON(fiber.Map{
			"short_url": "http://localhost:3000/" + code,
		})
	})

	// redirect
	app.Get("/:code", func(c *fiber.Ctx) error {
		code := c.Params("code")

		var longURL string

		err := db.QueryRow(
			"SELECT long_url FROM urls WHERE short_code=$1",
			code,
		).Scan(&longURL)

		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "URL not found",
			})
		}

		return c.Redirect(longURL, 302)
	})

	app.Listen(":3000")
}