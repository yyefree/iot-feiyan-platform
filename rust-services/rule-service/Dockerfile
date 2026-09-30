FROM docker.1ms.run/library/rust:1.75-alpine AS builder

WORKDIR /app
COPY Cargo.toml Cargo.lock ./
COPY common/ common/
COPY rule-service/ rule-service/

RUN cargo build --release --package iot-rule-service

FROM docker.1ms.run/library/alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/target/release/iot-rule-service .
EXPOSE 8083
ENTRYPOINT ["./iot-rule-service"]
