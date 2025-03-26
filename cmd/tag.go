/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strings"

	"github.com/adpg0222/aws-k8s/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: retag,
}

func retag(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		log.Fatalln("Please provide 2 arguments: source and target; repository must be the same")
	}

	if !strings.Contains(args[0], ":") || !strings.Contains(args[1], ":") {
		log.Fatalln("The arguments are incorrect; format: REPOSITORY:TAG")
	}

	client := ecr.NewFromConfig(aws.GetAwsConfig())
	ecr := aws.EcrService{Client: client}

	source := strings.Split(string(args[0]), ":")
	target := strings.Split(string(args[1]), ":")

	if source[0] != target[0] {
		log.Fatalln("The source and target repository are not the same")
	}

	imageManifest := ecr.GetImageManifest(source[0], source[1])
	_, err := ecr.PutImage(target[0], target[1], imageManifest)
	if err != nil {
		log.Fatalf("Failed to retag %s -> %s", args[0], args[1])
	}
	log.Printf("Tag created %s", args[1])
}

func init() {
	rootCmd.AddCommand(tagCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tagCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tagCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
