FROM rust:1.85-alpine AS builder

WORKDIR /app
COPY common/ common/
COPY ops-service/ ops-service/

WORKDIR /app/ops-service
RUN cargo build --release

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/ops-service/target/release/iot-ops-service .
EXPOSE 8085
ENTRYPOINT ["./iot-ops-service"]
