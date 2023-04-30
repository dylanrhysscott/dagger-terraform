# Dagger Terraform

A POC project experimenting with Dagger [https://dagger.io/](https://dagger.io/). This project creates a CLI which executes a Terraform pipleine as part of CI CD using pipelines as code.

## Requirements

* Docker
* Go

## Running the PoC


* Pipeline exec
    * From the `terraform-ci-runner` directory run `go build`
    * Execute `./terraform-ci-runner run` this will run Terraform init plan and apply
* Building binary as a Docker image using Dagger
    * From `container` folder run `go run .`. This will package the binary up as a Docker container and push it to Dockerhub
