package client

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

type Client struct {
	client api.Client
	api    v1.API
	cfg    *config.Config
}

func New(cfg *config.Config) (*Client, error) {
	client, err := api.NewClient(api.Config{
		Address: cfg.PrometheusUrl,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to initialize Prometheus client: %w", err)
	}

	return &Client{
		client: client,
		api:    v1.NewAPI(client),
		cfg:    cfg,
	}, nil
}

func (c *Client) Query(ctx context.Context, query string, time time.Time) (model.Vector, error) {
	res, wrn, err := c.api.Query(ctx, query, time)
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

func (c *Client) QueryRange(ctx context.Context, query string, start, end time.Time) (model.Matrix, error) {
	res, wrn, err := c.api.QueryRange(ctx, query, v1.Range{
		Start: start,
		End:   end,
		Step:  time.Second * time.Duration(c.cfg.QueryRangeStepSeconds),
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
