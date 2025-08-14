package utils

import (
	"context"
	"log/slog"
	"time"

	"connectrpc.com/connect"
)

// a generic helper to create a streaming RPC for a specific statistic.
func Stream[TRes any](
	ctx context.Context,
	stream *connect.ServerStream[TRes],
	fetchFn func(context.Context) (*TRes, error),
	intervalSeconds int,
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

	ticker := time.NewTicker(time.Duration(intervalSeconds) * time.Second)
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
