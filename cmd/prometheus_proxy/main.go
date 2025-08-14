package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"spitikos/api/internal/config"
	"spitikos/api/internal/logger"
	"spitikos/api/internal/prometheus_proxy/server"

	"buf.build/gen/go/spitikos/api/connectrpc/go/prometheus_proxy/v1/prometheus_proxyv1connect"
	"connectrpc.com/grpcreflect"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	logger.Init()

	cfg, err := config.New()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	mux := http.NewServeMux()

	// Add reflection
	reflector := grpcreflect.NewStaticReflector(prometheus_proxyv1connect.PrometheusServiceName)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	s, err := server.New(cfg)
	if err != nil {
		slog.Error("failed to create gRPC server", slog.Any("error", err))
		os.Exit(1)
	}
	path, handler := prometheus_proxyv1connect.NewPrometheusServiceHandler(s)
	mux.Handle(path, handler)

	// Add CORS middleware
	c := cors.New(cors.Options{
		AllowedMethods: []string{http.MethodGet, http.MethodPost, "OPTIONS"},
		AllowedHeaders: []string{"*"},
		AllowedOrigins: []string{"*"},
	})
	handlerWithCors := c.Handler(h2c.NewHandler(mux, &http2.Server{}))

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Port)
	slog.Info("server starting", "address", addr)

	if err := http.ListenAndServe(addr, handlerWithCors); err != nil {
		slog.Error("failed to listen and serve", slog.Any("error", err))
	}
}
