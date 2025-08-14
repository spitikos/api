# --- Deps Stage ---
# This stage is dedicated to downloading dependencies and is cached separately.
FROM golang:1.24-alpine AS deps
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download

FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY --from=deps /go/pkg/mod /go/pkg/mod
COPY . .
RUN go build -o /bin/server ./cmd/server

FROM gcr.io/distroless/static-debian12 AS runner
COPY --from=builder /bin/server /
EXPOSE 50051
CMD ["/server"]
