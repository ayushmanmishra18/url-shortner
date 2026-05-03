# 🔗 URL Shortener (Go + Fiber + PostgreSQL)

A high-performance URL shortener backend built using **Golang**, **Fiber**, and **PostgreSQL (Supabase)**.
Designed with production-level practices including **database-driven ID generation**, **Base62 encoding**, and **secure environment configuration**.

---

## 🚀 Features

* 🔗 Shorten long URLs into compact links
* ⚡ Fast redirection using optimized queries
* 🧠 Base62 encoding for clean, scalable short codes
* 🗄️ Persistent storage using PostgreSQL (Supabase)
* 🔒 Secure configuration using environment variables
* 🧪 REST API with proper validation & error handling
* ♻️ Restart-safe (no in-memory dependency)

---

## 🏗️ Architecture

```text
Client → Fiber (Go Backend) → PostgreSQL (Supabase)
```

---

## 🔁 How It Works

### 1. Shorten URL

* Client sends a long URL
* Server validates input
* Stores URL in database
* Retrieves auto-generated ID
* Converts ID → Base62 short code
* Updates DB with short code
* Returns short URL

---

### 2. Redirect

* User accesses short URL
* Server extracts short code
* Fetches original URL from DB
* Redirects user (HTTP 302)

---

## 🧠 Tech Stack

* **Backend:** Golang
* **Framework:** Fiber
* **Database:** PostgreSQL (Supabase)
* **Driver:** lib/pq
* **Env Management:** godotenv

---

## ⚙️ Setup & Installation

### 1. Clone the repository

```bash
git clone https://github.com/your-username/url-shortener.git
cd url-shortener
```

---

### 2. Install dependencies

```bash
go mod tidy
```

---

### 3. Create `.env` file

```env
DATABASE_URL=your_postgres_connection_string
```

---

### 4. Run the server

```bash
go run main.go
```

Server will start at:

```
http://localhost:3000
```

---

## 📡 API Endpoints

### 🔹 Health Check

```http
GET /health
```

---

### 🔹 Shorten URL

```http
POST /shorten
Content-Type: application/json
```

#### Request

```json
{
  "url": "https://example.com"
}
```

#### Response

```json
{
  "short_url": "http://localhost:3000/abc"
}
```

---

### 🔹 Redirect

```http
GET /:code
```

---

## 🗄️ Database Schema

```sql
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_code TEXT UNIQUE,
    long_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 🔐 Environment Variables

| Variable     | Description                  |
| ------------ | ---------------------------- |
| DATABASE_URL | PostgreSQL connection string |

---

## ⚠️ Current Limitations

* No caching (Redis planned)
* No rate limiting
* No analytics (click tracking)
* No custom aliases
* No expiration support

---

## 🚀 Future Improvements

* ⚡ Redis caching layer (performance boost)
* 📊 Click analytics & tracking
* 🔐 Authentication & user-based links
* 🌐 Custom domains
* ⏳ Link expiration support
* 🧱 Microservices architecture

---

## 🧪 Example Workflow

```text
POST /shorten → DB insert → ID → Base62 → store → return short URL
GET /abc → DB lookup → redirect to original URL
```

---

## 📌 Key Learning Outcomes

* Backend development in Go
* REST API design using Fiber
* PostgreSQL integration (cloud DB)
* Environment-based configuration
* Debugging real-world backend issues
* System design fundamentals

---

## 👨‍💻 Author

**Ayushman Mishra**
CSE @ JIIT | Backend & Cloud Enthusiast

---

## ⭐ Contribute / Feedback

Feel free to fork, improve, or raise issues.
If you found this helpful, consider giving a ⭐

---
