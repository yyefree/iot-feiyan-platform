FROM rust:1.85-alpine AS builder

WORKDIR /app
COPY common/ common/
COPY gateway/ gateway/

WORKDIR /app/gateway
RUN cargo build --release

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/gateway/target/release/iot-gateway .
EXPOSE 8080
ENTRYPOINT ["./iot-gateway"]
