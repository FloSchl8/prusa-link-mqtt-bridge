# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o prusa-link-mqtt-bridge ./cmd/prusa-link-mqtt-bridge

# Final stage
FROM scratch

COPY --from=builder /app/prusa-link-mqtt-bridge .

ENTRYPOINT ["./prusa-link-mqtt-bridge"]
