/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-ini/ini"
	"github.com/spf13/cobra"
)

const awsCredentialsFile string = "~/.aws/credentials"

// awsProfileCmd represents the awsProfile command
var awsProfileCmd = &cobra.Command{
	Use:     "awsProfile",
	Short:   "Select an AWS profile",
	Long:    "Select an AWS profile from you local credentials file.",
	Aliases: []string{"sp"},
	Run:     selectProfile,
}

func selectProfile(cmd *cobra.Command, args []string) {
	profiles := retrieveProfiles()

	var qs = []*survey.Question{
		{
			Name: "Profile",
			Prompt: &survey.Select{
				Default: "default",
				Message: "Choose a profile:",
				Options: profiles,
			},
		},
	}

	answers := struct {
		Profile string `survey:"Profile"`
	}{}
	err := survey.Ask(qs, &answers)
	if err != nil {
		if err.Error() == "interrupt" {
			log.Fatalf("ℹ Alright then, keep your profiles!\n")
		} else {
			log.Fatal(err.Error())
		}
	}
	log.Printf("export AWS_PROFILE=%s", answers.Profile)
}

func retrieveProfiles() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	awsCredentialsFile := path.Join(home, ".aws/credentials")

	credFile, err := ini.Load(awsCredentialsFile)
	if err != nil {
		log.Fatalf("❌ Failed to load AWS config file %s", awsCredPath)
	}

	sections := []string{}
	for _, section := range credFile.Sections() {
		sections = append(sections, section.Name())
	}
	return sections
}

func init() {
	rootCmd.AddCommand(awsProfileCmd)
}
