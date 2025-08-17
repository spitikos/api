package hello

import (
	"context"
	"spitikos/api/internal/config"
	"spitikos/api/internal/utils"

	hellov1 "buf.build/gen/go/spitikos/api/protocolbuffers/go/hello/v1"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	cfg *config.Config
}

func New(cfg *config.Config) (*Service, error) {
	return &Service{cfg}, nil
}

func (s *Service) Hello(
	ctx context.Context,
	req *connect.Request[hellov1.HelloRequest],
) (*connect.Response[hellov1.HelloResponse], error) {
	res := connect.NewResponse(&hellov1.HelloResponse{
		Reply: "Hello!",
	})
	return res, nil
}

func (s *Service) MyNameIs(
	ctx context.Context,
	req *connect.Request[hellov1.MyNameIsRequest],
) (*connect.Response[hellov1.MyNameIsResponse], error) {
	res := connect.NewResponse(&hellov1.MyNameIsResponse{
		Reply: "Your name is " + req.Msg.Name,
	})
	return res, nil
}

func (s *Service) Time(
	ctx context.Context,
	req *connect.Request[hellov1.TimeRequest],
	stream *connect.ServerStream[hellov1.TimeResponse],
) error {
	fetchFn := func(ctx context.Context) (*hellov1.TimeResponse, error) {
		res := hellov1.TimeResponse_builder{
			Now: timestamppb.Now(),
		}
		return res.Build(), nil
	}
	return utils.Stream(ctx, s.cfg, stream, fetchFn)
}
