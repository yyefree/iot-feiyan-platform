FROM docker.1ms.run/library/rust:1.75-alpine AS builder

WORKDIR /app
COPY Cargo.toml Cargo.lock ./
COPY common/ common/
COPY ops-service/ ops-service/

RUN cargo build --release --package iot-ops-service

FROM docker.1ms.run/library/alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/target/release/iot-ops-service .
EXPOSE 8085
ENTRYPOINT ["./iot-ops-service"]
