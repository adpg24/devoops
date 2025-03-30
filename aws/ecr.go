package aws

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

type EcrService struct {
	Client *ecr.Client
}

func (s *EcrService) GetImageManifest(repository string, imageTag string) string {
	image, err := s.GetImage(repository, imageTag)
	if err != nil {
		log.Fatalf("Failed to retrieve image '%s:%s': %v", repository, imageTag, err)
	}
	return *image.ImageManifest
}

func (s *EcrService) GetImage(repository string, imageTag string) (*types.Image, error) {
	input := ecr.BatchGetImageInput{ImageIds: []types.ImageIdentifier{{ImageTag: &imageTag}}, RepositoryName: &repository}
	output, err := s.Client.BatchGetImage(context.Background(), &input)
	if err != nil {
		return nil, err
	}
	return &output.Images[0], err
}

func (s *EcrService) PutImage(repository string, imageTag string, imageManifest string) (*types.Image, error) {
	input := ecr.PutImageInput{RepositoryName: &repository, ImageTag: &imageTag, ImageManifest: &imageManifest}
	output, err := s.Client.PutImage(context.Background(), &input)
	if err != nil {
		return nil, err
	}
	log.Println()
	return output.Image, nil
}
