FROM mcr.microsoft.com/dotnet/core/sdk:3.1.102-alpine as builder
WORKDIR /app
COPY . .
RUN dotnet restore && \
    dotnet build && \
    dotnet publish -c release -r linux-musl-x64 -o /cartservice

# cartservice
FROM alpine:3.10.4

RUN GRPC_HEALTH_PROBE_VERSION=v0.2.0 && \
    wget -qO/bin/grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /bin/grpc_health_probe

# Dependencies for runtime
# busybox-extras => telnet
RUN apk add --no-cache \
    busybox-extras \
    libunwind \
    libuuid \
    libgcc \
    libstdc++ \
    libintl \
    icu \
    libressl-dev
RUN echo 'http://dl-cdn.alpinelinux.org/alpine/v3.8/main' >> /etc/apk/repositories && apk update --no-cache && apk add --no-cache bash libc6-compat=1.1.19-r11

WORKDIR /app
COPY --from=builder /cartservice .
ENTRYPOINT ["./cartservice", "start"]
