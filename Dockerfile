FROM golang:1.24-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o /bin/server ./cmd/server

FROM gcr.io/distroless/static-debian12 AS runner
COPY --from=builder /bin/server /

# Expose the port the server listens on. This is for documentation purposes.
EXPOSE 50051

# Set the default command to run when the container starts.
CMD ["/server"]
