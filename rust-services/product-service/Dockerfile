FROM docker.1ms.run/library/rust:1.75-alpine AS builder

WORKDIR /app
COPY Cargo.toml Cargo.lock ./
COPY common/ common/
COPY product-service/ product-service/

RUN cargo build --release --package iot-product-service

FROM docker.1ms.run/library/alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/target/release/iot-product-service .
EXPOSE 8082
ENTRYPOINT ["./iot-product-service"]
