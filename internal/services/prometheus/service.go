package prometheus

import (
	"context"
	"fmt"
	"spitikos/api/internal/config"
	"spitikos/api/internal/utils"
	"time"

	v1 "buf.build/gen/go/spitikos/api/protocolbuffers/go/prometheus/v1"
	"connectrpc.com/connect"
)

type Service struct {
	cfg    *config.Config
	client *Client
}

func New(cfg *config.Config) (*Service, error) {
	client, err := NewClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus client: %w", err)
	}
	return &Service{cfg, client}, nil
}

func (s *Service) Query(
	ctx context.Context,
	req *connect.Request[v1.QueryRequest],
) (*connect.Response[v1.QueryResponse], error) {
	vector, err := s.client.Query(ctx, req.Msg.Query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to run Prometheus query: %w", err)
	}

	res := connect.NewResponse(VectorToQueryResponse(vector))
	return res, nil
}

func (s *Service) QueryRange(
	ctx context.Context,
	req *connect.Request[v1.QueryRangeRequest],
) (*connect.Response[v1.QueryRangeResponse], error) {
	matrix, err := s.client.QueryRange(ctx, req.Msg.Query, req.Msg.Since.AsTime())
	if err != nil {
		return nil, fmt.Errorf("failed to run Prometheus query range: %w", err)
	}
	res := connect.NewResponse(MatrixToQueryRangeResponse(matrix))
	return res, nil
}

func (s *Service) StreamQuery(
	ctx context.Context,
	req *connect.Request[v1.QueryRequest],
	stream *connect.ServerStream[v1.QueryResponse],
) error {
	fetchFn := func(ctx context.Context) (*v1.QueryResponse, error) {
		vector, err := s.client.Query(ctx, req.Msg.Query, time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to run Prometheus query: %w", err)
		}
		return VectorToQueryResponse(vector), nil
	}

	return utils.Stream(ctx, s.cfg, stream, fetchFn)
}

func (s *Service) StreamQueryRange(
	ctx context.Context,
	req *connect.Request[v1.QueryRangeRequest],
	stream *connect.ServerStream[v1.QueryRangeResponse],
) error {
	fetchFn := func(ctx context.Context) (*v1.QueryRangeResponse, error) {
		matrix, err := s.client.QueryRange(ctx, req.Msg.Query, req.Msg.Since.AsTime())
		if err != nil {
			return nil, fmt.Errorf("failed to run Prometheus query range: %w", err)
		}
		return MatrixToQueryRangeResponse(matrix), nil
	}

	return utils.Stream(ctx, s.cfg, stream, fetchFn)
}
