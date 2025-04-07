package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/ecr/types"
)

type EcrService struct {
	Client *ecr.Client
}

func (s *EcrService) GetImageManifest(repository string, imageTag string) (string, error) {
	image, err := s.GetImage(repository, imageTag)
	if err != nil {
		return "", err
	}
	return *image.ImageManifest, nil
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
	return output.Image, nil
}
