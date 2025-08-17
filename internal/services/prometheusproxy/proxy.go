package prometheusproxy

import (
	"context"
	"fmt"
	"log/slog"
	"spitikos/api/internal/config"
	"time"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

type Proxy struct {
	client api.Client
	api    v1.API
	cfg    *config.Config
}

func NewProxy(cfg *config.Config) (*Proxy, error) {
	client, err := api.NewClient(api.Config{
		Address: cfg.PrometheusProxy.Url,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize Prometheus client: %w", err)
	}

	return &Proxy{
		client: client,
		api:    v1.NewAPI(client),
		cfg:    cfg,
	}, nil
}

func (p *Proxy) Query(ctx context.Context, query string, time time.Time) (model.Vector, error) {
	res, wrn, err := p.api.Query(ctx, query, time)
	if err != nil {
		return nil, err
	}
	if len(wrn) > 0 {
		slog.Warn("Prometheus query completed with warnings", slog.Any("warnings", wrn))
	}

	vector, ok := res.(model.Vector)
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}

	return vector, nil
}

func (p *Proxy) QueryRange(ctx context.Context, query string, since time.Time) (model.Matrix, error) {
	res, wrn, err := p.api.QueryRange(ctx, query, v1.Range{
		Start: since,
		End:   time.Now(),
		Step:  time.Second * time.Duration(p.cfg.PrometheusProxy.QueryRangeStepSeconds),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to run Prometheus query range: %w", err)
	}
	if len(wrn) > 0 {
		slog.Warn("Prometheus query range completed with warnings", slog.Any("warnings", wrn))
	}

	matrix, ok := res.(model.Matrix)
	if !ok {
		return nil, fmt.Errorf("unexpected response type: %T", res)
	}

	return matrix, nil
}
