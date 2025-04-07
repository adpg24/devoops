/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/adpg24/devoops/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/spf13/cobra"
)

var repository string

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag [sourceTag] [newTag]",
	Short: "Retag an AWS ECR image",
	Long:  `Retag an AWS ECR image - No downloads necessary!`,
	Run:   run,
	Args:  cobra.MinimumNArgs(2),
}

func run(cmd *cobra.Command, args []string) {
	cfg, err := aws.GetAwsConfig(&aws.AwsConfig{})
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}

	client := ecr.NewFromConfig(*cfg)
	_ecr := aws.EcrService{Client: client}
	imageManifest, err := _ecr.GetImageManifest(repository, args[0])
	if err != nil {
		log.Fatalf("Failed to retrieve image manifest for %s:%s", repository, args[0])
	}

	_, err = _ecr.PutImage(repository, args[1], imageManifest)
	if err != nil {
		log.Fatalf("Failed to retag %s:%s -> %s:%s", repository, args[0], repository, args[1])
	}
	log.Printf("Tag created %s:%s", repository, args[1])
}

func init() {
	rootCmd.AddCommand(tagCmd)
	tagCmd.PersistentFlags().StringVarP(&repository, "repository", "r", "", "The ECR repository name")
}
