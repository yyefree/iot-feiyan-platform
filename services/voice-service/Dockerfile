# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o service .

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache curl

WORKDIR /app

COPY --from=builder /app/service .

EXPOSE 8089

CMD ["./service"]
