package internal

import (
	"context"
	"daggertf/ci"
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

func (t *TerraformCIRunner) RunPipeline(ctx context.Context, pipeline string) error {
	var client *dagger.Client
	var err error
	client, err = ci.CreateDaggerClient(ctx)
	if err != nil {
		return err
	}
	t.daggerClient = client
	switch pipeline {
	case "plan":
		err = t.createPlanPipeline(ctx)
	case "deploy":
		err = t.createDeployPipeline(ctx)
	default:
		err = fmt.Errorf("unknown pipeline: %s", pipeline)
	}
	return err
}

func (t *TerraformCIRunner) createPlanPipeline(ctx context.Context) error {
	pipeline := t.daggerClient.Pipeline("plan", dagger.PipelineOpts{
		Description: "A pipleine for planning Terraform",
	})
	defer pipeline.Close()
	var container *dagger.Container
	var err error
	container = t.terraformContainer(pipeline)
	container, err = ci.CreatePipelineStep(ctx, container, []string{"init"})
	if err != nil {
		return err
	}
	_, err = ci.CreatePipelineStep(ctx, container, []string{"plan", "-out", "server.plan"})
	if err != nil {
		return err
	}
	return nil
}

func (t *TerraformCIRunner) createDeployPipeline(ctx context.Context) error {
	pipeline := t.daggerClient.Pipeline("deploy", dagger.PipelineOpts{
		Description: "A pipleine for deploying Terraform",
	})
	defer pipeline.Close()
	var container *dagger.Container
	var err error
	container = t.terraformContainer(pipeline)
	container, err = ci.CreatePipelineStep(ctx, container, []string{"init"})
	if err != nil {
		return err
	}
	container, err = ci.CreatePipelineStep(ctx, container, []string{"plan", "-out", "server.plan"})
	if err != nil {
		return err
	}
	_, err = ci.CreatePipelineStep(ctx, container, []string{"apply", "server.plan"})
	if err != nil {
		return err
	}
	return nil
}

func (t *TerraformCIRunner) terraformContainer(client *dagger.Client) *dagger.Container {
	workingDirPath := "/src"
	fullImageTag := fmt.Sprintf("hashicorp/terraform:%s", t.ImageTag)
	sourceFolder := client.Host().Directory(t.SourceDirectory)
	return client.
		Container().
		From(fullImageTag).
		WithDirectory(workingDirPath, sourceFolder).
		WithWorkdir(workingDirPath)
}
