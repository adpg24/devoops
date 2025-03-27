package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

func GetAwsConfig() aws.Config {
	conf, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		log.Fatalln("Failed to load config: %v", err)
	}
	return conf
}
