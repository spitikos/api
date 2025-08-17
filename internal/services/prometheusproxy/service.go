package prometheusproxy

import (
	"context"
	"fmt"
	"spitikos/api/internal/config"
	"spitikos/api/internal/utils"
	"time"

	prometheusproxyv1 "buf.build/gen/go/spitikos/api/protocolbuffers/go/prometheusproxy/v1"
	"connectrpc.com/connect"
)

type Service struct {
	cfg   *config.Config
	proxy *Proxy
}

func New(cfg *config.Config) (*Service, error) {
	proxy, err := NewProxy(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus client: %w", err)
	}
	return &Service{cfg, proxy}, nil
}

func (s *Service) Query(
	ctx context.Context,
	req *connect.Request[prometheusproxyv1.QueryRequest],
	stream *connect.ServerStream[prometheusproxyv1.QueryResponse],
) error {
	fetchFn := func(ctx context.Context) (*prometheusproxyv1.QueryResponse, error) {
		vector, err := s.proxy.Query(ctx, req.Msg.Query, time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to run Prometheus query: %w", err)
		}
		return VectorToQueryResponse(vector), nil
	}

	return utils.Stream(ctx, s.cfg, stream, fetchFn)
}

func (s *Service) QueryRange(
	ctx context.Context,
	req *connect.Request[prometheusproxyv1.QueryRangeRequest],
	stream *connect.ServerStream[prometheusproxyv1.QueryRangeResponse],
) error {
	fetchFn := func(ctx context.Context) (*prometheusproxyv1.QueryRangeResponse, error) {
		matrix, err := s.proxy.QueryRange(ctx, req.Msg.Query, req.Msg.Since.AsTime())
		if err != nil {
			return nil, fmt.Errorf("failed to run Prometheus query range: %w", err)
		}
		return MatrixToQueryRangeResponse(matrix), nil
	}

	return utils.Stream(ctx, s.cfg, stream, fetchFn)
}
