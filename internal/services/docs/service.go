package docs

import (
	"context"
	"spitikos/api/internal/config"

	docspb "buf.build/gen/go/spitikos/api/protocolbuffers/go/docs"
	"connectrpc.com/connect"
)

type Service struct {
	cfg *config.Config
}

func New(cfg *config.Config) (*Service, error) {
	return &Service{cfg}, nil
}

func (s *Service) Sync(
	ctx context.Context,
	req *connect.Request[docspb.SyncRequest],
) (*connect.Response[docspb.SyncResponse], error) {
	builder := docspb.SyncResponse_builder{
		Success: true,
		Message: "Hello from Docs Service!",
	}
	res := connect.NewResponse(builder.Build())
	return res, nil
}

func (s *Service) GetSlugs(
	ctx context.Context,
	req *connect.Request[docspb.GetSlugsRequest],
) (*connect.Response[docspb.GetSlugsResponse], error) {
	builder := docspb.GetSlugsResponse_builder{
		Slugs: []string{"example-slug-1", "example-slug-2"},
	}
	res := connect.NewResponse(builder.Build())
	return res, nil
}
