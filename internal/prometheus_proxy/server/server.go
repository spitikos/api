package server

import (
	"context"
	"fmt"
	"spitikos/api/internal/config"
	"spitikos/api/internal/prometheus_proxy/client"
	"spitikos/api/internal/utils"
	"time"

	"buf.build/gen/go/spitikos/api/connectrpc/go/prometheus_proxy/v1/prometheus_proxyv1connect"
	pb "buf.build/gen/go/spitikos/api/protocolbuffers/go/prometheus_proxy/v1"
	"connectrpc.com/connect"
)

type Server struct {
	prometheus_proxyv1connect.UnimplementedPrometheusServiceHandler
	client *client.Client
	cfg    *config.PrometheusProxyConfig
}

func New(cfg *config.Config) (*Server, error) {
	client, err := client.New(&cfg.PrometheusProxy)
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus client: %w", err)
	}
	return &Server{
		client: client,
		cfg:    &cfg.PrometheusProxy,
	}, nil
}

func (s *Server) Query(
	ctx context.Context,
	req *connect.Request[pb.QueryRequest],
	stream *connect.ServerStream[pb.QueryResponse],
) error {
	fetchFn := func(ctx context.Context) (*pb.QueryResponse, error) {
		vector, err := s.client.Query(ctx, req.Msg.Query, time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to run Prometheus query: %w", err)
		}
		return VectorToQueryResponse(vector), nil
	}

	return utils.Stream(ctx, stream, fetchFn, s.cfg.StreamIntervalSeconds)
}

func (s *Server) QueryRange(
	ctx context.Context,
	req *connect.Request[pb.QueryRangeRequest],
	stream *connect.ServerStream[pb.QueryRangeResponse],
) error {
	fetchFn := func(ctx context.Context) (*pb.QueryRangeResponse, error) {
		matrix, err := s.client.QueryRange(ctx, req.Msg.Query, req.Msg.Since.AsTime())
		if err != nil {
			return nil, fmt.Errorf("failed to run Prometheus query range: %w", err)
		}
		return MatrixToQueryRangeResponse(matrix), nil
	}

	return utils.Stream(ctx, stream, fetchFn, s.cfg.StreamIntervalSeconds)
}
