package main

import (
	"context"
	"daggertf/ci"
	"log"

	"dagger.io/dagger"
)

func BuildActionImage() {
	ctx := context.TODO()
	log.Println("Creating Dagger Client...")
	client, err := ci.CreateDaggerClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	build(ctx, client)

}

func build(ctx context.Context, client *dagger.Client) {
	sourceFolder := client.Host().Directory(".")
	workingDirPath := "/src"
	log.Println("Building terraform-ci-runner...")
	build := client.Pipeline("action-build").
		Container().
		From("golang:latest").
		WithDirectory(workingDirPath, sourceFolder).
		WithWorkdir("/src/terraform-ci-runner").
		WithEnvVariable("CGO_ENABLED", "0").
		WithExec([]string{"go", "build", "-o", "terraform-ci-runner"})

	action := client.Pipeline("action-finalise").
		Container().
		From("alpine:latest").
		WithFile("/terraform-ci-runner", build.File("/src/terraform-ci-runner/terraform-ci-runner")).
		WithEntrypoint([]string{"/terraform-ci-runner", "--help"})

	addr, err := action.Publish(ctx, "dylanrhysscott/dagger-terraform:latest")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(addr)
}
