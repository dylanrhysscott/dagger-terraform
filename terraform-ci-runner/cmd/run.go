/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"daggertf/terraform-ci-runner/internal"
	"log"

	"github.com/spf13/cobra"
)

var version string
var sourceDir string
var pipeline string

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Runs a Terraform deployment",
	Long:  `Runs a Terraform deployment`,
	Run: func(cmd *cobra.Command, args []string) {
		runner, err := internal.NewTerraformCIRunner(context.TODO(), version, sourceDir)
		if err != nil {
			log.Fatal(err)
		}
		err = runner.RunPipeline(context.TODO(), pipeline)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVarP(&pipeline, "plan", "p", "deploy", "Pipeline to run")
	runCmd.Flags().StringVarP(&version, "source", "s", ".", "Path to the Terraform source")
	runCmd.Flags().StringVarP(&version, "version", "v", "1.4.6", "The image tag for hashicorp/terraform to use")
}
