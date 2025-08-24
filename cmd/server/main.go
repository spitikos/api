package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"spitikos/api/internal/config"
	"spitikos/api/internal/logger"
	"spitikos/api/internal/services/hello"
	"spitikos/api/internal/services/prometheusproxy"

	"buf.build/gen/go/spitikos/api/connectrpc/go/hello/v1/hellov1connect"
	"buf.build/gen/go/spitikos/api/connectrpc/go/prometheusproxy/v1/prometheusproxyv1connect"
	connectcors "connectrpc.com/cors"
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

	reflector := grpcreflect.NewStaticReflector(
		hellov1connect.HelloServiceName,
		prometheusproxyv1connect.PrometheusProxyServiceName,
	)

	helloSvc, err := hello.New(cfg)
	if err != nil {
		slog.Error("failed to create Hello server", slog.Any("error", err))
		os.Exit(1)
	}
	prometheusProxySvc, err := prometheusproxy.New(cfg)
	if err != nil {
		slog.Error("failed to create Prometheus Proxy server", slog.Any("error", err))
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	mux.Handle(hellov1connect.NewHelloServiceHandler(helloSvc))
	mux.Handle(prometheusproxyv1connect.NewPrometheusProxyServiceHandler(prometheusProxySvc))

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port)

	handler := h2c.NewHandler(mux, &http2.Server{})
	handler = withCORS(handler)

	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	slog.Info("server starting", slog.String("address", addr))
	if err := s.ListenAndServe(); err != nil {
		slog.Error("failed to listen and serve", slog.Any("error", err))
	}
}

func withCORS(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: connectcors.AllowedMethods(),
		AllowedHeaders: connectcors.AllowedHeaders(),
		ExposedHeaders: connectcors.ExposedHeaders(),
	})
	return c.Handler(h)
}
