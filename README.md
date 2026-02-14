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
POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_USER=postgres
POSTGRES_PASSWORD=example
POSTGRES_DB=postgres
REPOSITORY_TYPE=POSTGRES
HANDLER_TYPE=HTTP
PAGINATION_DEFAULT_LIMIT=20
PAGINATION_MAX_LIMIT=100
```

### 3. Run the application

```bash
go run ./cmd
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
curl http://localhost:8080/stock?page=1&limit=10
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
curl http://localhost:8080/restock/priorities?page=1&limit=10
```

---

## Running Tests

```bash
go test ./...
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
swag init -g cmd/main.go -o docs --parseInternal
```
