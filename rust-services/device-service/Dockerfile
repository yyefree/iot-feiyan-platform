FROM docker.1ms.run/library/rust:1.75-alpine AS builder

WORKDIR /app
ENV CARGO_REGISTRIES_CRATES_IO_PROTOCOL=sparse
ENV CARGO_TARGET_X86_64_UNKNOWN_LINUX_GNU_LINKER=x86_64-linux-gnu-gcc

# Install build dependencies
RUN apk add --no-cache musl-dev openssl-dev openssl-libs-static

COPY Cargo.toml Cargo.lock ./
COPY common/ common/
COPY device-service/ device-service/

RUN cargo build --release --package iot-device-service

FROM docker.1ms.run/library/alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Shanghai
WORKDIR /app

COPY --from=builder /app/target/release/iot-device-service .

EXPOSE 8081
ENTRYPOINT ["./iot-device-service"]
