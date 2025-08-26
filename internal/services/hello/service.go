package hello

import (
	"context"
	"spitikos/api/internal/config"
	"spitikos/api/internal/utils"

	hellopb "buf.build/gen/go/spitikos/api/protocolbuffers/go/hello"
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
	req *connect.Request[hellopb.HelloRequest],
) (*connect.Response[hellopb.HelloResponse], error) {
	builder := hellopb.HelloResponse_builder{
		Reply: "Hello!",
	}
	res := connect.NewResponse(builder.Build())
	return res, nil
}

func (s *Service) MyNameIs(
	ctx context.Context,
	req *connect.Request[hellopb.MyNameIsRequest],
) (*connect.Response[hellopb.MyNameIsResponse], error) {
	builder := hellopb.MyNameIsResponse_builder{
		Reply: "Your name is " + req.Msg.Name,
	}
	res := connect.NewResponse(builder.Build())
	return res, nil
}

func (s *Service) Time(
	ctx context.Context,
	req *connect.Request[hellopb.TimeRequest],
	stream *connect.ServerStream[hellopb.TimeResponse],
) error {
	fetchFn := func(ctx context.Context) (*hellopb.TimeResponse, error) {
		res := hellopb.TimeResponse_builder{
			Now: timestamppb.Now(),
		}
		return res.Build(), nil
	}
	return utils.Stream(ctx, s.cfg, stream, fetchFn)
}
