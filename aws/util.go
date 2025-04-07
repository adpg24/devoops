package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type AwsConfig struct {
	Region  string
	Profile string
}

func GetAwsConfig(awsConfig *AwsConfig) (*aws.Config, error) {
	var optFns []func(*config.LoadOptions) error
	if awsConfig.Profile != "" {
		optFns = append(optFns, config.WithSharedConfigProfile(awsConfig.Profile))
	}

	if awsConfig.Region != "" {
		optFns = append(optFns, config.WithRegion(awsConfig.Region))
	}

	conf, err := config.LoadDefaultConfig(context.TODO(), optFns...)
	if err != nil {
		return nil, err
	}

	return &conf, nil
}
