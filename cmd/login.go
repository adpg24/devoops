/*
Copyright © 2024 Antonio Pizarro adpg0222@gmail.com
*/
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"time"

	"github.com/adpg24/devoops/aws"
	"github.com/adpg24/devoops/util"
	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-ini/ini"
)

type mfaSurveyAnswer = struct {
	MfaDevice string `survey:"mfaDevice"`
	MfaCode   string `survey:"mfaCode"`
}

var (
	awsCredPath      string
	awsProfile       string
	longTermProfile  string
	shortTermProfile string
	region           = "eu-west-1"
	mfaDevice        string
)

const (
	longTermSuffix        string = "-mfa"
	keyAwsAccessKey       string = "aws_access_key_id"
	keyAwsSecretAccessKey string = "aws_secret_access_key"
)

// loginCmd represents the login command
var loginCmd = &cobra.Command{
	Use:     "login",
	Aliases: []string{"log", "l"},
	Short:   "Generate short time credentials with MFA authentication",
	Run:     login,
	PreRun:  checkFlags,
}

func checkFlags(cmd *cobra.Command, args []string) {
	if awsProfile == "default" {
		if awsProfileEnv := os.Getenv("AWS_PROFILE"); awsProfileEnv != "" {
			awsProfile = awsProfileEnv
		}
	}
}

func login(cmd *cobra.Command, args []string) {
	// load INI files (~/.aws/credentials)
	credFile, err := ini.Load(awsCredPath)
	util.HandleErr(err, "❌ Failed to load AWS config file %s", awsCredPath)

	shortTermProfile = awsProfile
	longTermProfile = fmt.Sprintf("%s%s", awsProfile, longTermSuffix)

	// validate long term profile = [profile]-mfa
	if longTermCreds, err := credFile.GetSection(longTermProfile); err != nil {
		log.Fatalf("❌ AWS Profile not available! Please create a long-term profile with the suffix \"-mfa\". e.g. [default] -> [default-mfa]\n")
	} else {
		requiredKeys := []string{keyAwsAccessKey, keyAwsSecretAccessKey}
		for _, key := range requiredKeys {
			if !longTermCreds.HasKey(key) {
				log.Fatalf("❌ The profile %s does not have the key '%s'\n", longTermProfile, key)
			}
		}

		if configRegion, err := longTermCreds.GetKey("region"); err == nil {
			region = configRegion.String()
		}
		if configMfaDevice, err := longTermCreds.GetKey("aws_mfa_device"); err == nil {
			mfaDevice = configMfaDevice.String()
		}
	}

	// validate short term profile = [profile]
	if shortTermConfig, err := credFile.GetSection(shortTermProfile); err == nil {
		currentTime := time.Now()
		expirationkey := shortTermConfig.Key("expiration")

		if expiration, err := time.Parse("2006-01-02 15:04:05", expirationkey.String()); err != nil {
			log.Fatalf("❌ Expiration (%s) in profile \"%s\" is in the wrong format (2006-01-02 15:04:05)!\nError: %s\n", expirationkey.String(), shortTermProfile, err.Error())
		} else if expiration.After(currentTime) {
			log.Printf("ℹ You're still authenticated! Your credential will expire at %s.\n", expiration.Format("2006-01-02 15:04:05"))
			os.Exit(0)
		}
	}

	conf, err := aws.GetAwsConfig(&aws.AwsConfig{Region: region, Profile: longTermProfile})
	util.HandleErr(err, "Failed to retrieve config: %v", err)

	var qs = []*survey.Question{
		{
			Name:     "mfaCode",
			Prompt:   &survey.Input{Message: "Please enter the MFA code for the given MFA device:"},
			Validate: survey.ComposeValidators(survey.MinLength(6), survey.MaxLength(6), survey.Required),
		},
	}

	var answers mfaSurveyAnswer

	if mfaDevice == "" {
		_iam := iam.NewFromConfig(*conf)

		devices, err := _iam.ListMFADevices(context.TODO(), &iam.ListMFADevicesInput{})
		util.HandleErr(err, "❌ An error occurred while listing MFA devices: %v", err)

		var mfaDevices []string

		if len(devices.MFADevices) < 1 {
			log.Fatalf("❌ No mfa devices have been configured for this user!")
		}

		for _, device := range devices.MFADevices {
			mfaDevices = append(mfaDevices, *device.SerialNumber)
		}

		q := &survey.Question{
			Name: "mfaDevice",
			Prompt: &survey.Select{
				Message: "Choose a MFA device:",
				Options: mfaDevices,
			},
		}

		// insert(qs, q, 0)
		qs = append(qs[:0], append([]*survey.Question{q}, qs[0:]...)...)
		answers = mfaSurveyAnswer{}
	} else {
		answers = mfaSurveyAnswer{MfaDevice: mfaDevice}
	}

	err = survey.Ask(qs, &answers)
	if err != nil {
		if err.Error() == "interrupt" {
			log.Fatalf("ℹ Alright then, keep your secrets! Exiting..\n")
		} else {
			log.Fatal(err.Error())
		}
	}

	_sts := sts.NewFromConfig(*conf)
	session, err := _sts.GetSessionToken(context.TODO(), &sts.GetSessionTokenInput{
		TokenCode:    &answers.MfaCode,
		SerialNumber: &answers.MfaDevice,
	})
	util.HandleErr(err, "❌ An error occurred while retrieving the session token for %s!: %v", longTermProfile, err)

	// func Section will create the profile (INI section) if it does not exist
	sectionKeys := map[string]string{
		"aws_access_key_id":     *session.Credentials.AccessKeyId,
		"aws_secret_access_key": *session.Credentials.SecretAccessKey,
		"aws_session_token":     *session.Credentials.SessionToken,
		"expiration":            session.Credentials.Expiration.Format("2006-01-02 15:04:05"),
	}
	err = util.AddProfileSection(awsCredPath, credFile, shortTermProfile, sectionKeys)
	util.HandleErr(err, "Failed to add new profile/section to %s: %v", awsCredPath, err)

	log.Printf("The short-term credentials were successfully created for profile %s", shortTermProfile)
}

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// log.SetPrefix("devoops\t")
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	home, err := os.UserHomeDir()
	util.HandleErr(err, "Failed to retrieve use home dir: %v", err)

	loginCmd.Flags().StringVarP(&awsCredPath, "config", "c", path.Join(home, ".aws/credentials"), "AWS credentials file location")
	loginCmd.Flags().StringVarP(&awsProfile, "profile", "p", "default", "AWS profile for which you need to authenticate with MFA")
}
