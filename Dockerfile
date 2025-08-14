FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /bin/prometheus-proxy ./cmd/prometheus_proxy

FROM gcr.io/distroless/static-debian12 as runner
COPY --from=builder /bin/prometheus-proxy /
