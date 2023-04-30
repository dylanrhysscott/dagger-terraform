package ci

import (
	"context"
	"fmt"

	"dagger.io/dagger"
)

type TerraformCIRunner struct {
	daggerClient    *dagger.Client
	ImageTag        string
	SourceDirectory string
}

func NewTerraformCIRunner(ctx context.Context, terraformImageTag string, sourceDir string) (*TerraformCIRunner, error) {
	return &TerraformCIRunner{
		ImageTag:        terraformImageTag,
		SourceDirectory: sourceDir,
	}, nil
}

func (t *TerraformCIRunner) RunPipeline(ctx context.Context) error {
	client, err := createDaggerClient(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	t.daggerClient = client
	c := t.terraformContainer()
	c, err = t.runInitStep(ctx, c)
	if err != nil {
		return err
	}
	c, err = t.runPlanStep(ctx, c)
	if err != nil {
		return err
	}
	_, err = t.runApplyStep(ctx, c)
	if err != nil {
		return err
	}
	return nil
}

func (t *TerraformCIRunner) runInitStep(ctx context.Context, c *dagger.Container) (*dagger.Container, error) {
	c = c.WithExec([]string{"init"})
	out, err := c.Stdout(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println(out)
	return c, nil
}

func (t *TerraformCIRunner) runPlanStep(ctx context.Context, c *dagger.Container) (*dagger.Container, error) {
	c = c.WithExec([]string{"plan", "-out", "server.plan"})
	out, err := c.Stdout(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println(out)
	return c, nil
}

func (t *TerraformCIRunner) runApplyStep(ctx context.Context, c *dagger.Container) (*dagger.Container, error) {
	c = c.WithExec([]string{"apply", "server.plan"})
	out, err := c.Stdout(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println(out)
	return c, nil
}

func (t *TerraformCIRunner) terraformContainer() *dagger.Container {
	workingDirPath := "/src"
	fullImageTag := fmt.Sprintf("hashicorp/terraform:%s", t.ImageTag)
	sourceFolder := t.daggerClient.Host().Directory(t.SourceDirectory)
	return t.daggerClient.
		Container().
		From(fullImageTag).
		WithDirectory(workingDirPath, sourceFolder).
		WithWorkdir(workingDirPath)
}
