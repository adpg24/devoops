package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var activeProfile = "defaul"

type EcsTool struct {
	client  *ecs.ECS
	profile string
	output  string
	region  string
}

func getActiveProfile() string {
	activeProfile = os.Getenv("AWS_PROFILE")
	return activeProfile
}

func (config *EcsTool) init() {
	mySession := session.Must(session.NewSession())
	config.client = ecs.New(mySession)
}

func (tool *EcsTool) listCluster() []*string {
	maxRes := int64(5)
	input := ecs.ListClustersInput{MaxResults: &maxRes}
	response, err := tool.client.ListClusters(&input)
	if err != nil {
		log.Fatalf("Failed to retrieve clusters")
	}
	clusterArns := response.ClusterArns
	return clusterArns
}

//Hello from this fihhjjjje
