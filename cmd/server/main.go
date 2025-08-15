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

	// Register the reflection service.
	// This reflector needs to know about all services that will be registered.
	reflector := grpcreflect.NewStaticReflector(
		prometheus_proxyv1connect.PrometheusServiceName,
		// Add new service names here in the future.
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))

	// Register the Prometheus Proxy service.
	prometheusProxyServer, err := server.New(cfg)
	if err != nil {
		slog.Error("failed to create prometheus proxy server", slog.Any("error", err))
		os.Exit(1)
	}
	path, handler := prometheus_proxyv1connect.NewPrometheusServiceHandler(prometheusProxyServer)
	mux.Handle(path, handler)

	// handlerWithCors := c.Handler(h2c.NewHandler(mux, &http2.Server{}))

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port)
	slog.Info("server starting", "address", addr)

	if err := http.ListenAndServe(addr, h2c.NewHandler(mux, &http2.Server{})); err != nil {
		slog.Error("failed to listen and serve", slog.Any("error", err))
	}
}
