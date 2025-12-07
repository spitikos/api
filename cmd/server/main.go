package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"spitikos/api/internal/config"
	"spitikos/api/internal/services/docs"
	"spitikos/api/internal/services/hello"
	"spitikos/api/internal/services/prometheus"
	"time"

	"buf.build/gen/go/spitikos/api/connectrpc/go/docs/docsconnect"
	"buf.build/gen/go/spitikos/api/connectrpc/go/hello/helloconnect"
	"buf.build/gen/go/spitikos/api/connectrpc/go/prometheus/prometheusconnect"
	"connectrpc.com/grpcreflect"
	"github.com/lmittmann/tint"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("failed to load config", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("config loaded", slog.Any("config", cfg))

	helloSvc, err := hello.New(cfg)
	if err != nil {
		slog.Error("failed to create Hello server", slog.Any("error", err))
		os.Exit(1)
	}
	prometheusSvc, err := prometheus.New(cfg)
	if err != nil {
		slog.Error("failed to create Prometheus server", slog.Any("error", err))
		os.Exit(1)
	}
	docsSvc, err := docs.New(cfg)
	if err != nil {
		slog.Error("failed to create Docs server", slog.Any("error", err))
		os.Exit(1)
	}

	slog.Info("services initialized")

	reflector := grpcreflect.NewStaticReflector(
		helloconnect.HelloServiceName,
		prometheusconnect.PrometheusServiceName,
		docsconnect.DocsServiceName,
	)

	mux := http.NewServeMux()
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	mux.Handle(helloconnect.NewHelloServiceHandler(helloSvc))
	mux.Handle(prometheusconnect.NewPrometheusServiceHandler(prometheusSvc))
	mux.Handle(docsconnect.NewDocsServiceHandler(docsSvc))

	addr := fmt.Sprintf("0.0.0.0:%d", cfg.Server.Port)
	handler := h2c.NewHandler(mux, &http2.Server{})

	s := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	slog.Info("server starting", slog.String("address", addr))

	if err := s.ListenAndServe(); err != nil {
		slog.Error("failed to listen and serve", slog.Any("error", err))
	}
}

func init() {
	w := os.Stderr

	slog.SetDefault(slog.New(
		tint.NewHandler(w, &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.Kitchen,
		}),
	))
}
