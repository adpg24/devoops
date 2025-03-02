package main

import (
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var activeProfile = "defaul"

type EcsService struct {
	client  *ecs.ECS
	profile string
	output  string
	region  string
}

func getActiveProfile() string {
	activeProfile = os.Getenv("AWS_PROFILE")
	return activeProfile
}

func (config *EcsService) init() {
	mySession := session.Must(session.NewSession())
	config.client = ecs.New(mySession)
}

func (tool *EcsService) listClusters() []*string {
	maxRes := int64(5)
	input := ecs.ListClustersInput{MaxResults: &maxRes}
	response, err := tool.client.ListClusters(&input)
	if err != nil {
		log.Fatalf("Failed to retrieve clusters")
	}
	return response.ClusterArns
}

func (tool *EcsService) listServices(cluster string) []*string {
	input := ecs.ListServicesInput{Cluster: &cluster}
	response, err := tool.client.ListServices(&input)
	if err != nil {
		log.Fatalf("Failed to retrieve clusters")
	}
	return response.ServiceArns
}

func (tool *EcsService) describeService(cluster string, service string) *ecs.Service {
	services := []*string{&service}
	input := ecs.DescribeServicesInput{Cluster: &cluster, Services: services}
	response, err := tool.client.DescribeServices(&input)
	if err != nil {
		log.Fatalf("Failed to retrieve clusters")
	}
	return response.Services[0]
}
