# Firelink

A Go-powered REST API for recipes, drinks, books, and more—built with Gin, GORM, and PostgreSQL.
Includes notification support, healthchecks, and robust Docker/CI integration.

---

## Description

**Firelink** is a modular, extensible API server for food, drink, and book discovery.
It features endpoints for random dinner recipes (via Spoonacular), cocktail recipes, book lookups (Gutenberg), and more.
The project is production-ready, with healthchecks, Prometheus metrics, notification support (ntfy), and a clean, testable codebase.

---

## Features

- 🍽️ **Dinner Recipes:** Get random or specific recipes from Spoonacular.
- 🍸 **Bartender:** Random cocktail recipes, save to DB, and view history.
- 📚 **Books:** Check for books in the Gutenberg project.
- 🩺 **Healthcheck:** Simple endpoint for monitoring.
- 📝 **Notifications:** Send rich notifications via ntfy.
- 🐳 **Dockerized:** Easy to run locally or in production.
- 🧪 **CI/CD:** GitHub Actions for build, test, and coverage.

---

## Getting Started

### Dependencies

- Go 1.23+
- Docker & Docker Compose (for local development)
- PostgreSQL (default: runs in Docker)
- [Spoonacular API Key](https://spoonacular.com/food-api)

### Installing

1. **Clone the repo:**
   ```sh
   git clone https://github.com/rjhoppe/firelink.git
   cd firelink
   ```

2. **Set up environment variables:**
   - Copy `.env.example` to `.env` and fill in your values:
     ```
     POSTGRES_USER=youruser
     POSTGRES_PASSWORD=yourpassword
     POSTGRES_DB=firelink
     SPOONACULAR_API_KEY=your_spoonacular_key
     ```

3. **Start with Docker Compose:**
   ```sh
   docker-compose up --build
   ```

4. **Or run locally:**
   ```sh
   go run main.go
   ```

### Executing program

- **API available at:** `http://localhost:8080`
- **Healthcheck:** `GET /healthcheck`
- **Help/Docs:** `GET /help`

---

## API Endpoints

See `/help` endpoint for a full, live list.
Example endpoints:

- `GET /dinner/random` — 3 random dinner recipes
- `GET /dinner/recipe/:id` — Recipe by ID
- `GET /bartender/random` — Random cocktail
- `POST /bartender/save` — Save last cocktail to DB
- `GET /bartender/history` — Cocktail history
- `GET /ebook/find/:title` — Check for a book
- `POST /database/backup` — Backup the database

---

## Testing

Run all tests with coverage:
```sh
go test -v -cover ./...
```

---

## Help

- For troubleshooting, see logs or run:
  ```sh
  docker-compose logs
  ```
- For environment variable issues, check `.env` and your Docker Compose config.

---