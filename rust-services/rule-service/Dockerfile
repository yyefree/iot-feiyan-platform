FROM rust:1.85-alpine AS builder

WORKDIR /app
COPY common/ common/
COPY rule-service/ rule-service/

WORKDIR /app/rule-service
RUN cargo build --release

FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /app/rule-service/target/release/iot-rule-service .
EXPOSE 8083
ENTRYPOINT ["./iot-rule-service"]
