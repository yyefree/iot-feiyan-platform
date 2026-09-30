FROM docker.1ms.run/library/golang:1.23-alpine AS builder

WORKDIR /app
ENV GOPROXY=https://goproxy.cn,direct
ENV GONOSUMCHECK=*
ENV GONOSUMDB=*
ENV GOFLAGS=-mod=mod
COPY go.mod ./
RUN GOPROXY=https://goproxy.cn,direct GONOSUMCHECK=* GONOSUMDB=* go mod download || true
COPY . .
RUN GOPROXY=https://goproxy.cn,direct GONOSUMCHECK=* GONOSUMDB=* go build -mod=mod -ldflags="-s -w" -o /app/api-gateway

FROM docker.1ms.run/library/alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/api-gateway .
EXPOSE 8080
ENTRYPOINT ["./api-gateway"]
