# URL Shortener Backend

A high-performance, normalized URL shortening service built with Go, PostgreSQL, Redis, and Consul.

## 🚀 Features
- **Normalized Schema**: Efficient storage using 5 tables (`users`, `long_urls`, `short_urls`, `clicks`, `url_stats`).
- **High Performance Redirection**: Cache-aside strategy with Redis for O(1) lookups.
- **Async Analytics**: Non-blocking click tracking using a worker pool and Go channels.
- **Config Management**: Dynamic configuration via Consul and Viper.
- **OAuth2 Authentication**: Secure registration and token-based access.

## 🛠️ Prerequisites
- Docker (for PostgreSQL, Redis, and Consul)
- Go 1.21+

## ⚙️ Initial Setup

### 1. Infrastructure
```bash
docker start pg-container redis-stack consul
```

### 2. Configuration Seeding
```bash
bash scripts/seed_consul.sh
```

### 3. Database Migrations
You can run migrations using the built-in Go CLI or via a bash command for manual schema management.

**Using Go CLI (Recommended):**
```bash
export CONSUL_URL=localhost:8500
export CONSUL_PATH=config/backend
go run main.go migrate
```

**Using Bash (Manual Rollback/Setup):**
```bash
# To drop all normalized tables
docker exec -i pg-container psql -U postgres -d url_shortener -c "DROP TABLE IF EXISTS url_stats, clicks, short_urls, long_urls, users CASCADE;"
```

## 🏃 Running the Application
```bash
./run.sh
```

---

## 📚 Documentation
- **[API Testing Guide](api_testing.md)**: Real-world CURL commands for all 13 endpoints.
- **[Performance Monitoring](performance_monitoring.md)**: How to track latency and worker pool health.
- **[Operations Guide](operations_guide.md)**: Detailed setup and deployment instructions.
