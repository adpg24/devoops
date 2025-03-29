/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"slices"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/adpg24/devoops/util"
	"github.com/go-ini/ini"
	"github.com/spf13/cobra"
)

const awsCredentialsFile string = ".aws/credentials"

type AwsProfile struct {
	Name    string
	Account string
}

func (p *AwsProfile) String() string {
	return fmt.Sprintf("%s/%s", p.Account, p.Name)
}

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

	profileOptions := []string{}
	for _, p := range profiles {
		profileOptions = append(profileOptions, p.String())
	}
	sort.Strings(profileOptions)

	var qs = []*survey.Question{
		{
			Name: "Profile",
			Prompt: &survey.Select{
				Default: profileOptions[0],
				Message: "Choose a profile:",
				Options: profileOptions,
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
	selectedProfile := strings.Split(answers.Profile, "/")[1]
	exportCmd := fmt.Sprintf("export AWS_PROFILE=%s", selectedProfile)
	err = util.CopyToClipboard(exportCmd)
	if err != nil {
		log.Fatalln("Somethin went wrong while copying to clipboard:", err)
	}
	log.Println("Export command written to clipboard")
}

func retrieveProfiles() []AwsProfile {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	awsCredentialsFile := path.Join(home, awsCredentialsFile)

	credFile, err := ini.Load(awsCredentialsFile)
	if err != nil {
		log.Fatalf("❌ Failed to load AWS config file %s", awsCredPath)
	}

	profiles := []AwsProfile{}

	for _, section := range credFile.Sections() {
		accountId := "????????????"
		var keyWithAccountId string
		if slices.Contains(section.KeyStrings(), "role_arn") {
			keyWithAccountId = "role_arn"
		} else if slices.Contains(section.KeyStrings(), "aws_mfa_device") {
			keyWithAccountId = "aws_mfa_device"
		}

		if keyWithAccountId != "" {
			key, _ := section.GetKey(keyWithAccountId)
			if key != nil {
				accountId = strings.Split(key.String(), ":")[4]
			}
		}
		profiles = append(profiles, AwsProfile{Name: section.Name(), Account: accountId})
	}
	return profiles
}

func init() {
	rootCmd.AddCommand(awsProfileCmd)
}
