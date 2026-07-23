FROM golang:1.25 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -v -o /usr/local/bin/app ./cmd

FROM golang:1.25

COPY --from=builder /usr/local/bin/app /usr/local/bin/app
COPY --from=builder /usr/src/app/internal/infrastructure/repository/db/seed.sql /usr/local/bin/seed.sql

EXPOSE 8080

CMD ["app"]