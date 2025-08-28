package githubclient

import (
	"context"
	"errors"
	"fmt"
	"spitikos/api/internal/config"

	"github.com/google/go-github/v74/github"
)

type Client struct {
	client *github.Client
	cfg    *config.Config
}

func New(cfg *config.Config) (*Client, error) {
	client := github.NewClient(nil).WithAuthToken(cfg.Github.Token)

	return &Client{
		client: client,
		cfg:    cfg,
	}, nil
}

func (c *Client) GetContent(ctx context.Context, repo string, path string) (string, error) {
	file, _, res, err := c.client.Repositories.GetContents(ctx, c.cfg.Github.Owner, repo, path, nil)
	if err != nil {
		return "", err
	}
	if res.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch content of %s/%s", repo, path)
	}

	content, err := file.GetContent()
	if err != nil {
		return "", err
	}

	return content, nil
}

func (c *Client) GetCommit(ctx context.Context, repo string, sha string) (*github.Commit, error) {
	commit, _, err := c.client.Repositories.GetCommit(ctx, c.cfg.Github.Owner, repo, sha, nil)
	if err != nil {
		return nil, err
	}

	return commit.Commit, nil
}

func (c *Client) GetLatestCommit(ctx context.Context, repo string, path string) (*github.Commit, error) {
	commits, _, err := c.client.Repositories.ListCommits(ctx, c.cfg.Github.Owner, repo, &github.CommitsListOptions{
		Path: path,
	})
	if err != nil {
		return nil, err
	}
	if len(commits) == 0 {
		return nil, errors.New("no commits found")
	}

	return commits[0].Commit, nil
}

func (c *Client) GetRateLimit(ctx context.Context) (*github.Rate, error) {
	rateLimits, _, err := c.client.RateLimit.Get(ctx)
	if err != nil {
		return nil, err
	}

	return rateLimits.Core, nil
}
