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

	"github.com/spf13/cobra"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-ini/ini"
)

type mfaSurveyAnswer = struct {
	MfaDevice string `survey:"mfaDevice"`
	MfaCode   string `survey:"mfaCode"`
}

type awsConfig = struct {
	AwsAccessKey       string
	AwsSecretAccessKey string
	Region             string
	AwsSessionToken    string
	MfaDevice          string
	Expiration         string
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

func insert[T any](array []T, element T, i int) []T {
	return append(array[:i], append([]T{element}, array[i:]...)...)
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
	if err != nil {
		log.Fatalf("❌ Failed to load AWS config file %s", awsCredPath)
	}

	shortTermProfile = awsProfile
	longTermProfile = fmt.Sprintf("%s%s", awsProfile, longTermSuffix)

	// validate long term profile = [profile]-mfa
	if longTermCreds, err := credFile.GetSection(longTermProfile); err != nil {
		log.Fatalf("❌ AWS Profile not available! Please suffix the profile you want to use with \"-mfa\". e.g. [default] -> [default-mfa]\n")
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

	conf, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithSharedConfigProfile(longTermProfile),
	)
	if err != nil {
		log.Fatal(err)
	}

	_iam := iam.NewFromConfig(conf)

	devices, err := _iam.ListMFADevices(context.TODO(), &iam.ListMFADevicesInput{})
	if err != nil {
		log.Fatalf("❌ An error occurred while listing MFA devices!\nError: %s\n", err.Error())
	}

	var mfaDevices []string

	if len(devices.MFADevices) < 1 {
		log.Fatalf("❌ No mfa devices have been configured for this user!")
	}

	for _, device := range devices.MFADevices {
		mfaDevices = append(mfaDevices, *device.SerialNumber)
	}

	var qs = []*survey.Question{
		{
			Name:     "mfaCode",
			Prompt:   &survey.Input{Message: "Please enter the MFA code for the given MFA device:"},
			Validate: survey.ComposeValidators(survey.MinLength(6), survey.MaxLength(6), survey.Required),
		},
	}

	var answers mfaSurveyAnswer

	if mfaDevice == "" {
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

	_sts := sts.NewFromConfig(conf)
	session, err := _sts.GetSessionToken(context.TODO(), &sts.GetSessionTokenInput{
		TokenCode:    &answers.MfaCode,
		SerialNumber: &answers.MfaDevice,
	})
	if err != nil {
		log.Fatalf("❌ An error occurred while retrieving session token for %s!\nError: %s\n", answers.MfaDevice, err.Error())
		return
	}

	// func Section will create the profile (INI section) if it does not exist
	sec := credFile.Section(shortTermProfile)

	sec.NewKey("aws_access_key_id", *session.Credentials.AccessKeyId)
	sec.NewKey("aws_secret_access_key", *session.Credentials.SecretAccessKey)
	sec.NewKey("aws_session_token", *session.Credentials.SessionToken)
	sec.NewKey("expiration", session.Credentials.Expiration.Format("2006-01-02 15:04:05"))

	credFile.SaveTo(awsCredPath)
}

func SetVersionInfo(version, commit, date string) {
	rootCmd.Version = fmt.Sprintf("%s (Built on %s from Git SHA %s)", version, date, commit)
}

func init() {
	rootCmd.AddCommand(loginCmd)

	// log.SetPrefix("devoops\t")
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	// rootCmd.PersistentFlags().StringVar(&awsCredPath, "config", path.Join(home, ".aws/credentials"), "AWS credentials file location")
	loginCmd.Flags().StringVarP(&awsCredPath, "config", "c", path.Join(home, ".aws/credentials"), "AWS credentials file location")
	// rootCmd.PersistentFlags().StringVar(&awsProfile, "profile", "default", "AWS Profile for which we need to request a MFA token")
	loginCmd.Flags().StringVarP(&awsProfile, "profile", "p", "default", "AWS profile for which you need to authenticate with MFA")
}
