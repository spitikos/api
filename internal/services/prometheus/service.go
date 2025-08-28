package prometheus

import (
	"context"
	"fmt"
	prometheusclient "spitikos/api/internal/clients/prometheus"
	"spitikos/api/internal/config"
	"spitikos/api/internal/utils"
	"time"

	prometheuspb "buf.build/gen/go/spitikos/api/protocolbuffers/go/prometheus"
	"connectrpc.com/connect"
)

type Service struct {
	cfg        *config.Config
	prometheus *prometheusclient.Client
}

func New(cfg *config.Config) (*Service, error) {
	client, err := prometheusclient.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Prometheus client: %w", err)
	}
	return &Service{cfg, client}, nil
}

func (s *Service) Query(
	ctx context.Context,
	req *connect.Request[prometheuspb.QueryRequest],
) (*connect.Response[prometheuspb.QueryResponse], error) {
	vector, err := s.prometheus.Query(ctx, req.Msg.Query, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to run Prometheus query: %w", err)
	}

	res := connect.NewResponse(buildQueryResponse(vector))
	return res, nil
}

func (s *Service) QueryRange(
	ctx context.Context,
	req *connect.Request[prometheuspb.QueryRangeRequest],
) (*connect.Response[prometheuspb.QueryRangeResponse], error) {
	matrix, err := s.prometheus.QueryRange(ctx, req.Msg.Query, req.Msg.Since.AsTime())
	if err != nil {
		return nil, fmt.Errorf("failed to run Prometheus query range: %w", err)
	}
	res := connect.NewResponse(buildQueryRangeResponse(matrix))
	return res, nil
}

func (s *Service) StreamQuery(
	ctx context.Context,
	req *connect.Request[prometheuspb.StreamQueryRequest],
	stream *connect.ServerStream[prometheuspb.StreamQueryResponse],
) error {
	fetchFn := func(ctx context.Context) (*prometheuspb.StreamQueryResponse, error) {
		vector, err := s.prometheus.Query(ctx, req.Msg.Query, time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to run Prometheus query: %w", err)
		}
		res := buildStreamQueryResponse(vector)
		return res, nil
	}

	return utils.Stream(ctx, s.cfg, stream, fetchFn)
}

func (s *Service) StreamQueryRange(
	ctx context.Context,
	req *connect.Request[prometheuspb.StreamQueryRangeRequest],
	stream *connect.ServerStream[prometheuspb.StreamQueryRangeResponse],
) error {
	fetchFn := func(ctx context.Context) (*prometheuspb.StreamQueryRangeResponse, error) {
		matrix, err := s.prometheus.QueryRange(ctx, req.Msg.Query, req.Msg.Since.AsTime())
		if err != nil {
			return nil, fmt.Errorf("failed to run Prometheus query range: %w", err)
		}
		res := buildStreamQueryRangeResponse(matrix)
		return res, nil
	}

	return utils.Stream(ctx, s.cfg, stream, fetchFn)
}
