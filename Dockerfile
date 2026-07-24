FROM golang:1.25 AS builder

WORKDIR /usr/src/app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -v -o /usr/local/bin/app ./cmd/app
RUN CGO_ENABLED=0 go build -v -o /usr/local/bin/migrate ./cmd/migrate
RUN CGO_ENABLED=0 go build -v -o /usr/local/bin/seed ./cmd/seed

FROM scratch

COPY --from=builder /usr/local/bin/app /usr/local/bin/app
COPY --from=builder /usr/local/bin/migrate /usr/local/bin/migrate
COPY --from=builder /usr/local/bin/seed /usr/local/bin/seed

ENV PATH=/usr/local/bin

EXPOSE 8080

CMD ["app"]
