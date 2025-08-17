package utils

import (
	"context"
	"log/slog"
	"spitikos/api/internal/config"
	"time"

	"connectrpc.com/connect"
)

// a generic helper to create a streaming RPC for a specific statistic.
func Stream[TRes any](
	ctx context.Context,
	cfg *config.Config,
	stream *connect.ServerStream[TRes],
	fetchFn func(context.Context) (*TRes, error),
) error {
	// initial fetch
	data, err := fetchFn(ctx)
	if err != nil {
		slog.Error("failed to fetch", slog.Any("error", err))
		return err
	}
	if err := stream.Send(data); err != nil {
		slog.Error("failed to send", slog.Any("error", err))
		return err
	}

	ticker := time.NewTicker(time.Duration(cfg.Server.StreamIntervalSeconds) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			data, err := fetchFn(ctx)
			if err != nil {
				slog.Error("failed to fetch", slog.Any("error", err))
				return err
			}
			if err := stream.Send(data); err != nil {
				slog.Error("failed to send", slog.Any("error", err))
				return err
			}
		case <-ctx.Done():
			slog.Info("client disconnected")
			return nil
		}
	}
}
