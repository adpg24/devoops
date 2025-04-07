package aws

import (
	"os"
	"slices"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
)

func TestConfigRegion(t *testing.T) {
	// create a dummy credentials file
	// the credentials file contains one profile: test
	os.Setenv("AWS_SHARED_CREDENTIALS_FILE", "/tmp/credentials")

	credentialsContents := `
[test]
aws_access_key_id     = XXXXXX
aws_secret_access_key = XXXXXX
region                = eu-west-1
output                = json
aws_mfa_device        = arn:aws:iam::607344922194:mfa/test`

	err := os.WriteFile("/tmp/credentials", []byte(credentialsContents), 0644)
	if err != nil {
		t.Fatalf("Failed to write mock credentials")
	}

	inputRegion := "eu-west-2"
	inputProfile := "test"

	cfg, err := GetAwsConfig(&AwsConfig{Region: inputRegion, Profile: inputProfile})
	if err != nil {
		t.Fatalf("Failed with error %v", err)
	}

	if cfg.Region != inputRegion {
		t.Fatalf("The input region %s is not the same as the config.Region %s\n", inputRegion, cfg.Region)
	}

	var profiles []string
	for _, c := range cfg.ConfigSources {
		t, ok := c.(config.SharedConfig)
		if ok {
			profiles = append(profiles, t.Profile)
		}
	}

	if !slices.Contains(profiles, inputProfile) {
		t.Fatalf("Expected profile %s in the ConfigSources, but was not found", inputProfile)
	}
}

func TestConfigLoadFails(t *testing.T) {
	defer func() { _ = recover() }()

	inputProfile := "profileNotExists"
	if _, err := GetAwsConfig(&AwsConfig{Profile: inputProfile}); err == nil {
		t.Errorf("did not panic")
	}

}

func TestConfigEmptyStruct(t *testing.T) {
	_, err := GetAwsConfig(&AwsConfig{})
	if err != nil {
		t.Fatalf("Failed to get AWS config: %v", err)
	}
}
