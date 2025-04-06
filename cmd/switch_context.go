/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/adpg24/devoops/kube"
	"log"

	"github.com/spf13/cobra"
)

// switchContextCmd represents the switchContext command
var switchContextCmd = &cobra.Command{
	GroupID: "contextGroup",
	Use:     "switchContext",
	Aliases: []string{"sc"},
	Short:   "Switch context",
	Long:    `Switch to another context - the contexts are retrieved from you kube config`,
	Run:     runConfig,
}

func init() {
	rootCmd.AddCommand(switchContextCmd)
}

func runConfig(cmd *cobra.Command, args []string) {
	kubeConfig := kube.NewKubeConfig("")
	contexts := kubeConfig.GetContexts()

	var qs = []*survey.Question{
		{
			Name: "context",
			Prompt: &survey.Select{
				Default: kubeConfig.GetCurrentContext(),
				Message: "Choose a context:",
				Options: contexts,
			},
		},
	}

	answers := struct {
		Context string `survey:"Context"`
	}{}
	err := survey.Ask(qs, &answers)
	if err != nil {
		if err.Error() == "interrupt" {
			log.Fatalf("ℹ Alright then, keep your contexts!\n")
		} else {
			log.Fatal(err.Error())
		}
	}
	log.Printf("Switched to context %s", answers.Context)
	kubeConfig.SetCurrentContext(answers.Context)
}
