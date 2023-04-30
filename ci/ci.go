package ci

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

type CIRunner interface {
	RunPipeline(ctx context.Context, pipeline string) error
}

func CreateDaggerClient(ctx context.Context) (*dagger.Client, error) {
	client, err := dagger.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func CreatePipelineStep(ctx context.Context, c *dagger.Container, args []string) (*dagger.Container, error) {
	c = c.WithExec(args)
	out, err := c.Stdout(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println(out)
	return c, nil
}
