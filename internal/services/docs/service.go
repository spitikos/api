package docs

import (
	"context"
	"fmt"
	"slices"
	"spitikos/api/internal/config"
	"sync"

	githubclient "spitikos/api/internal/clients/github"
	redisclient "spitikos/api/internal/clients/redis"

	docspb "buf.build/gen/go/spitikos/api/protocolbuffers/go/docs"
	"connectrpc.com/connect"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	cfg    *config.Config
	redis  *redisclient.Client
	github *githubclient.Client
}

const (
	githubRepo      = "docs"
	docsSlugsSetKey = "docs:slugs"
)

func New(cfg *config.Config) (*Service, error) {
	redis, err := redisclient.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create Redis client: %w", err)
	}

	github, err := githubclient.New(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create GitHub client: %w", err)
	}

	return &Service{
		cfg:    cfg,
		redis:  redis,
		github: github,
	}, nil
}

func (s *Service) Sync(
	ctx context.Context,
	req *connect.Request[docspb.SyncRequest],
) (*connect.Response[docspb.SyncResponse], error) {
	staleKeys, err := s.redis.GetKeys(ctx, "docs:*")
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get existing doc keys: %w", err))
	}

	paths := req.Msg.Paths
	statusMap := &sync.Map{}
	g, gCtx := errgroup.WithContext(ctx)

	for _, path := range paths {
		g.Go(func() error {
			status := &docspb.SyncStatus{}
			content, err := s.github.GetContent(gCtx, githubRepo, path)
			if err != nil {
				return fmt.Errorf("failed to get content for path %s: %w", path, err)
			}
			status.GetContent = true

			commit, err := s.github.GetLatestCommit(gCtx, githubRepo, path)
			if err != nil {
				return fmt.Errorf("failed to get latest commit for path %s: %w", path, err)
			}
			status.GetLatestCommit = true

			doc := &docspb.Doc_builder{
				Title:    getFirstH1(content),
				Modified: timestamppb.New(commit.Author.Date.Time),
				Content:  content,
			}

			data, err := protojson.Marshal(doc.Build())
			if err != nil {
				return fmt.Errorf("failed to marshal doc for path %s: %w", path, err)
			}
			status.MarshalDoc = true

			key := slugToKey(trimExtension(path))
			if err := s.redis.Set(gCtx, key, string(data)); err != nil {
				return fmt.Errorf("failed to cache doc for path %s: %w", path, err)
			}
			status.CacheDoc = true

			statusMap.Store(path, status)

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	for _, staleKey := range staleKeys {
		slug := keyToSlug(staleKey)
		if !slices.Contains(paths, slug+".md") {
			if err := s.redis.Del(ctx, staleKey); err != nil {
				return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to delete stale doc key %s: %w", staleKey, err))
			}
		}
	}

	status := make(map[string]*docspb.SyncStatus)
	statusMap.Range(func(key, value any) bool {
		status[key.(string)] = value.(*docspb.SyncStatus)
		return true
	})

	res := connect.NewResponse(&docspb.SyncResponse{
		Status: status,
	})
	return res, nil
}

func (s *Service) GetSlugs(
	ctx context.Context,
	req *connect.Request[docspb.GetSlugsRequest],
) (*connect.Response[docspb.GetSlugsResponse], error) {
	keys, err := s.redis.GetKeys(ctx, "docs:*")
	if err != nil {
		return nil, fmt.Errorf("failed to get docs keys: %w", err)
	}

	slugs := make([]string, 0, len(keys))
	for _, key := range keys {
		slugs = append(slugs, keyToSlug(key))
	}

	builder := docspb.GetSlugsResponse_builder{
		Slugs: slugs,
	}

	res := connect.NewResponse(builder.Build())
	return res, nil
}

func (s *Service) GetDoc(
	ctx context.Context,
	req *connect.Request[docspb.GetDocRequest],
) (*connect.Response[docspb.GetDocResponse], error) {
	slug := req.Msg.Slug
	key := slugToKey(slug)

	data, err := s.redis.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get doc %s from redis: %w", key, err)
	}

	var doc docspb.Doc
	if err := protojson.Unmarshal([]byte(data), &doc); err != nil {
		return nil, fmt.Errorf("failed to unmarshal doc %s: %w", key, err)
	}

	res := connect.NewResponse(&docspb.GetDocResponse{
		Doc: &doc,
	})

	return res, nil
}
