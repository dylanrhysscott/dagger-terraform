/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"daggertf/ci"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var version string
var sourceDir string
var pipeline string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "daggertf",
	Short: "A Terraform CI runner built with Dagger",
	Long:  `A Terraform CI runner built with Dagger`,
	Run: func(cmd *cobra.Command, args []string) {
		runner, err := ci.NewTerraformCIRunner(context.TODO(), version, sourceDir)
		if err != nil {
			log.Fatal(err)
		}
		err = runner.RunPipeline(context.TODO(), pipeline)
		if err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&pipeline, "plan", "p", "deploy", "Pipeline to run")
	rootCmd.Flags().StringVarP(&version, "source", "s", ".", "Path to the Terraform source")
	rootCmd.Flags().StringVarP(&version, "version", "v", "1.4.6", "The image tag for hashicorp/terraform to use")
}
