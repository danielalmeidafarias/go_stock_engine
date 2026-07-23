# Stock Engine API

Microservice for managing product stock and calculating restock priorities.

## Requirements

- **With Docker:** Docker and Docker Compose
- **Without Docker:** Go 1.25+, PostgreSQL 16+

---

## Running with Docker

```bash
docker compose up --build
```

The API will be available at `http://localhost:8080` and Swagger UI at `http://localhost:8080/swagger/index.html`.

To stop:

```bash
docker compose down
```

To stop and remove the database volume:

```bash
docker compose down -v
```

---

## Running without Docker

### 1. Set up PostgreSQL

Make sure you have a PostgreSQL instance running locally.

### 2. Configure environment variables

```bash
cp .env.example .env
```

Edit `.env` with your database credentials:

```dotenv
DATABASE_DRIVER=POSTGRES
DATABASE_URL=postgres://postgres:example@localhost:5432/postgres?sslmode=disable
HANDLER_TYPE=HTTP
PAGINATION_DEFAULT_LIMIT=20
PAGINATION_MAX_LIMIT=100
PRIORITY_NEGATIVE_STOCK_FACTOR=1.5
PRIORITY_LEAD_TIME_FACTOR=1.2
PRIORITY_ZERO_SALES_FACTOR=0.5
REQUEST_TIMEOUT_ENABLED=true
REQUEST_TIMEOUT=5s
```

### 3. Run migrations

The API does not change the schema during startup. Apply pending migrations before
running it:

```bash
go run ./cmd/migrate
```

With Docker Compose, migrations run automatically before the API starts. To run
them manually:

```bash
docker compose run --rm migrate
```

### 4. Seed the database (optional)

The API never seeds data during startup. To load the demonstration products into an empty database, run:

```bash
go run ./cmd/seed
```

With Docker Compose:

```bash
docker compose run --rm app seed
```

The command skips execution when product stock records already exist.

### 5. Run the application

```bash
go run ./cmd/app
```

The API will be available at `http://localhost:8080`.

---

## API Endpoints

| Method | Route                         | Description                     |
|--------|-------------------------------|---------------------------------|
| POST   | `/stock`                      | Create a product stock          |
| GET    | `/stock`                      | List all product stocks         |
| GET    | `/stock/:id`                  | Get a product stock by ID       |
| PUT    | `/stock/:id`                  | Update a product stock          |
| DELETE | `/stock/:id`                  | Delete a product stock          |
| GET    | `/stock/category/:category`   | List product stocks by category |
| GET    | `/restock/priorities`         | Get restock priorities          |
| GET    | `/swagger/index.html`               | Swagger UI                      |

---

## Request Examples

### Create a product stock

```bash
curl -X POST http://localhost:8080/stock \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Oil Filter X",
    "category": "engine",
    "current_stock": 15,
    "minimum_stock": 20,
    "average_daily_sales": 4,
    "lead_time_days": 5,
    "unit_cost": 18.50,
    "criticality_level": 3
  }'
```

### List all product stocks

```bash
curl "http://localhost:8080/stock?page=1&limit=10"
```

### Get a product stock by ID

```bash
curl http://localhost:8080/stock/{id}
```

### Update a product stock

```bash
curl -X PUT http://localhost:8080/stock/{id} \
  -H "Content-Type: application/json" \
  -d '{
    "current_stock": 25,
    "minimum_stock": 30
  }'
```

### Delete a product stock

```bash
curl -X DELETE http://localhost:8080/stock/{id}
```

### Get restock priorities

```bash
curl "http://localhost:8080/restock/priorities?page=1&limit=10"
```

> **Design decision:** Restock priorities are calculated and ordered in memory before pagination. This keeps the priority rules and configurable policy factors in the domain layer instead of duplicating business rules in database queries.

---

## Running Tests

### Domain Restock Priority tests

```bash
go test -v ./internal/domain
```

### Use cases unit tests

```bash
go test -v ./internal/application/tests
```

### HTTP benchmarks

Requires the application running (`docker compose up -d` or `go run ./cmd/app`):

```bash
go test ./tests/benchmark -run '^$' -bench . -benchmem
```

### E2E contract tests

Runs against an isolated PostgreSQL environment and removes it afterward:

```bash
make e2e
```

---

## Swagger

With the application running, access the interactive API documentation at:

```
http://localhost:8080/swagger/index.html
```

### Regenerating Swagger docs

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/app/main.go -o docs --parseInternal
```
