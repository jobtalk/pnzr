package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

/*
cred := credentials.NewSharedCredentials("", "default")
awsConfig := &aws.Config{
	Credentials: cred,
	Region:      aws.String("ap-northeast-1"),
}
*/

func RegisterTaskDefinition(awsConfig *aws.Config, registerTaskDefinitionInput *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
	svc := ecs.New(session.New(), awsConfig)

	return svc.RegisterTaskDefinition(registerTaskDefinitionInput)
}

func CreateService(awsConfig *aws.Config, createServiceInput *ecs.CreateServiceInput) (*ecs.CreateServiceOutput, error) {
	svc := ecs.New(session.New(), awsConfig)
	return svc.CreateService(createServiceInput)
}

func UpdateService(awsConfig *aws.Config, updateServiceInput *ecs.UpdateServiceInput) (*ecs.UpdateServiceOutput, error) {
	svc := ecs.New(session.New(), awsConfig)
	return svc.UpdateService(updateServiceInput)
}
