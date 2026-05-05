package main

import (
	"context"
	"database/sql"
	"log"
	"net/url"
	"os"
	"time"

	_ "github.com/lib/pq"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var db *sql.DB
var rdb *redis.Client
var ctx = context.Background()

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

	// 🔹 Load .env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	app := fiber.New()

	// 🔹 DATABASE CONNECTION
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL not set")
	}

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("DB OPEN ERROR:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("DB PING ERROR:", err)
	}

	log.Println("✅ Connected to PostgreSQL")

	// REDIS CONNECTION
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal("❌ Redis connection failed:", err)
	}

	log.Println(" Connected to Redis")

	// HEALTH CHECK
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "Server is running",
		})
	})

	// SHORTEN URL
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

		// Step 1: Insert and get ID
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

		// Step 2: Generate short code
		code := encodeBase62(id)

		// Step 3: Update short_code
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

		// 🔥 Store in Redis (cache it immediately)
		rdb.Set(ctx, code, body.URL, 10*time.Minute)

		return c.JSON(fiber.Map{
			"short_url": "http://localhost:3000/" + code,
		})
	})

	//REDIRECT (WITH CACHE)
	app.Get("/:code", func(c *fiber.Ctx) error {
		code := c.Params("code")

		// 1. Check Redis
		val, err := rdb.Get(ctx, code).Result()
		if err == nil {
			log.Println("⚡ Cache HIT")
			return c.Redirect(val, 302)
		}

		// 2. Fallback to DB
		var longURL string
		err = db.QueryRow(
			"SELECT long_url FROM urls WHERE short_code=$1",
			code,
		).Scan(&longURL)

		if err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "URL not found",
			})
		}

		// 3. Store in Redis
		rdb.Set(ctx, code, longURL, 10*time.Minute)

		log.Println("🐢 Cache MISS → stored in Redis")

		return c.Redirect(longURL, 302)
	})

	// START SERVER
	log.Fatal(app.Listen(":3000"))
}