package ci

import (
	"context"

	"dagger.io/dagger"
)

type CIRunner interface {
	RunPipeline(ctx context.Context) error
}

func createDaggerClient(ctx context.Context) (*dagger.Client, error) {
	client, err := dagger.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}
