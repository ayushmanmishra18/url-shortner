🚀 URL Shortener (Go + Fiber + PostgreSQL + Redis)

A high-performance URL shortener built with Go (Fiber), PostgreSQL, and Redis caching.
Designed with real backend engineering practices like caching, database indexing, and scalable architecture.

🌟 Features
🔗 Shorten long URLs into compact links
⚡ Fast redirection using Redis cache
🗄️ Persistent storage with PostgreSQL
🧠 Base62 encoding for unique short codes
⏱️ Cache with TTL for optimized performance
📦 Clean REST API design
🌱 Environment-based configuration
🏗️ Tech Stack
Backend Framework: Fiber
Language: Go
Database: PostgreSQL
Cache Layer: Redis
Driver: lib/pq
Config: godotenv
🧠 Architecture
Client
  ↓
Fiber Server
  ↓
Redis (Cache Layer)
  ↓ (miss)
PostgreSQL (Persistent Storage)
⚙️ Setup Instructions
1. Clone the repository
git clone https://github.com/your-username/url-shortner.git
cd url-shortner
2. Install dependencies
go mod tidy
3. Create .env file
DATABASE_URL=your_postgresql_connection_string
REDIS_ADDR=localhost:6379
4. Run Redis

Make sure Redis is running on port 6379

5. Run the server
go run main.go

Server will start at:

http://localhost:3000
📌 API Endpoints
🔹 Health Check
GET /health
🔹 Shorten URL
POST /shorten
Body:
{
  "url": "https://example.com"
}
Response:
{
  "short_url": "http://localhost:3000/abc123"
}
🔹 Redirect
GET /:code

Redirects to original URL.

⚡ Caching Strategy
Type: Distributed Cache
Strategy: Cache-Aside (Lazy Loading)
Tool: Redis
Flow:
1st Request → DB → Cache → Response
Next Requests → Cache → Response (FAST)
🧪 Example
curl -X POST http://localhost:3000/shorten \
-H "Content-Type: application/json" \
-d '{"url":"https://google.com"}'
📊 Current Capabilities
Handles fast read-heavy workloads
Reduces DB load using caching
Scalable backend structure
🚧 Future Improvements
📊 Click analytics (tracking visits)
🔐 Rate limiting (Redis-based)
⏳ URL expiration
✨ Custom short URLs
🌍 Deployment (Docker + Cloud)
📈 Monitoring & logging
👨‍💻 Author

Ayushman Mishra
CSE @ JIIT | Backend & Cloud Enthusiast