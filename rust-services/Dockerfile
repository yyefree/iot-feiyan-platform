FROM docker.1ms.run/library/rust:1.75-alpine AS builder

WORKDIR /app
ENV CARGO_REGISTRIES_CRATES_IO_PROTOCOL=sparse
ENV CARGO_TARGET_X86_64_UNKNOWN_LINUX_GNU_LINKER=x86_64-linux-gnu-gcc

RUN apk add --no-cache musl-dev openssl-dev openssl-libs-static

COPY . .

# Build all services from workspace
RUN cargo build --release --bin iot-device-service
RUN cargo build --release --bin iot-product-service
RUN cargo build --release --bin iot-rule-service
RUN cargo build --release --bin iot-data-service
RUN cargo build --release --bin iot-ops-service
RUN cargo build --release --bin iot-gateway

FROM docker.1ms.run/library/alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Shanghai
WORKDIR /app

COPY --from=builder /app/target/release/iot-device-service /app/device-service
COPY --from=builder /app/target/release/iot-product-service /app/product-service
COPY --from=builder /app/target/release/iot-rule-service /app/rule-service
COPY --from=builder /app/target/release/iot-data-service /app/data-service
COPY --from=builder /app/target/release/iot-ops-service /app/ops-service
COPY --from=builder /app/target/release/iot-gateway /app/gateway

EXPOSE 8080 8081 8082 8083 8084 8085
CMD ["/app/gateway"]
