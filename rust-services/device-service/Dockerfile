FROM rust:1.85-alpine AS builder

WORKDIR /app
ENV CARGO_REGISTRIES_CRATES_IO_PROTOCOL=sparse
ENV CARGO_TARGET_X86_64_UNKNOWN_LINUX_GNU_LINKER=x86_64-linux-gnu-gcc

RUN apk add --no-cache musl-dev openssl-dev openssl-libs-static

COPY common/ common/
COPY device-service/ device-service/

WORKDIR /app/device-service
RUN cargo build --release

FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Shanghai
WORKDIR /app

COPY --from=builder /app/device-service/target/release/iot-device-service .

EXPOSE 8081
ENTRYPOINT ["./iot-device-service"]
