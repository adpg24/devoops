/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strings"

	"github.com/adpg24/devoops/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "Retag an AWS ECR image",
	Long: `Retag an AWS ECR image - No downloads necessary!
	devoops tag IMAGE_REPOSITORY_NAME:TAG -> IMAGE_REPOSITORY_NAME:NEW_TAG`,
	Run: run,
}

func run(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		log.Fatalln("Please provide 2 arguments: source and target; repository must be the same")
	}

	if !strings.Contains(args[0], ":") || !strings.Contains(args[1], ":") {
		log.Fatalln("The arguments are incorrect; format: REPOSITORY:TAG")
	}

	cfg, err := aws.GetAwsConfig(&aws.AwsConfig{})
	if err != nil {
		log.Fatalf("Failed to load AWS config: %v", err)
	}
	client := ecr.NewFromConfig(*cfg)
	ecr := aws.EcrService{Client: client}

	source := strings.Split(string(args[0]), ":")
	target := strings.Split(string(args[1]), ":")

	if source[0] != target[0] {
		log.Fatalln("The source and target repository are not the same")
	}

	imageManifest := ecr.GetImageManifest(source[0], source[1])
	_, err = ecr.PutImage(target[0], target[1], imageManifest)
	if err != nil {
		log.Fatalf("Failed to retag %s -> %s", args[0], args[1])
	}
	log.Printf("Tag created %s", args[1])
}

func init() {
	rootCmd.AddCommand(tagCmd)
}
