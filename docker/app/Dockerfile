# syntax=docker/dockerfile:1.7
ARG GO_VERSION=1.25.1

FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /src
RUN apk add --no-cache build-base git

# 1) сначала только go.mod/go.sum — для кэша
COPY go.mod go.sum ./
RUN go mod download

# 2) затем весь код
COPY .. .

# 3) сборка под платформу билдера
ENV CGO_ENABLED=0
RUN go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/booker

# --- runtime ---
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates tzdata
COPY --from=builder /out/app /usr/local/bin/app

# конфиги монтируем томом; путь — через ENV
ENV CONFIG_PATH=/config/local.yaml
ENTRYPOINT ["/usr/local/bin/app"]