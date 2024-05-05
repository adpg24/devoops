/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"

	"github.com/AlecAivazis/survey/v2"
	"github.com/adpg0222/aws-k8s/kube"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: runConfig,
}

func runConfig(cmd *cobra.Command, args []string) {
	kube := kube.KubeConfig{}
	contexts := kube.GetContexts()

	var qs = []*survey.Question{
		{
			Name: "context",
			Prompt: &survey.Select{
				Message: "Choose a context:",
				Options: contexts,
			},
		},
	}

	answers := struct{ context string }{}
	err := survey.Ask(qs, &answers)
	if err != nil {
		if err.Error() == "interrupt" {
			log.Fatalf("ℹ Alright then, keep your contexts! Exiting..\n")
		} else {
			log.Fatal(err.Error())
		}
	}

}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
