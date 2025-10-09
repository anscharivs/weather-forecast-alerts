FROM golang:1.25 AS builder
WORKDIR /app
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o weather-forecast-alerts ./cmd/server

FROM debian:bullseye-slim
WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates

COPY --from=builder /app/weather-forecast-alerts .
CMD ["./weather-forecast-alerts"]
