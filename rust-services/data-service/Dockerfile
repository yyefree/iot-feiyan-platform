FROM rust:1.85-alpine AS builder

WORKDIR /app
COPY common/ common/
COPY data-service/ data-service/

WORKDIR /app/data-service
RUN cargo build --release

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/data-service/target/release/iot-data-service .
EXPOSE 8084
ENTRYPOINT ["./iot-data-service"]
