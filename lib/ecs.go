package lib

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

type Service struct {
	Name *string
}

type ECS struct {
	svc ecsiface.ECSAPI
}

func NewECS(awsConfig *aws.Config) *ECS {
	return &ECS{
		svc: ecs.New(session.New(), awsConfig),
	}
}

func (e *ECS) RegisterTaskDefinition(registerTaskDefinitionInput *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
	return e.svc.RegisterTaskDefinition(registerTaskDefinitionInput)
}

func (e *ECS) ListServices(params *ecs.ListServicesInput) (*ecs.ListServicesOutput, error) {
	var (
		ret     = &ecs.ListServicesOutput{}
		pageNum int
	)

	err := e.svc.ListServicesPages(params, func(page *ecs.ListServicesOutput, lastPage bool) bool {
		pageNum++
		ret.ServiceArns = append(ret.ServiceArns, page.ServiceArns...)
		return pageNum <= 1000
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func (e *ECS) CreateService(createServiceInput *ecs.CreateServiceInput) (*ecs.CreateServiceOutput, error) {
	return e.svc.CreateService(createServiceInput)
}

func (e *ECS) UpdateService(updateServiceInput *ecs.UpdateServiceInput) (*ecs.UpdateServiceOutput, error) {
	return e.svc.UpdateService(updateServiceInput)
}

func UpsertService(awsConfig *aws.Config, createServiceInput *ecs.CreateServiceInput) (interface{}, error) {
	ecsClient := NewECS(awsConfig)
	if len(createServiceInput.LoadBalancers) != 0 && targetGroupARN != nil {
		createServiceInput.LoadBalancers[0].TargetGroupArn = targetGroupARN
	}
	if createServiceInput.Cluster == nil {
		createServiceInput.Cluster = aws.String("default")
	}
	ok, err := IsExistService(awsConfig, *createServiceInput.Cluster, *createServiceInput.ServiceName)
	if err != nil {
		return nil, err
	}
	if !ok {
		return ecsClient.CreateService(createServiceInput)
	}

	updateServiceInput := &ecs.UpdateServiceInput{}
	updateServiceInput.Cluster = createServiceInput.Cluster
	updateServiceInput.DeploymentConfiguration = createServiceInput.DeploymentConfiguration
	updateServiceInput.DesiredCount = createServiceInput.DesiredCount
	updateServiceInput.Service = createServiceInput.ServiceName
	updateServiceInput.TaskDefinition = createServiceInput.TaskDefinition
	return ecsClient.UpdateService(updateServiceInput)
}

func IsExistService(awsConfig *aws.Config, clusetrName string, serviceName string) (bool, error) {
	ecsClient := NewECS(awsConfig)
	listInput := &ecs.ListServicesInput{
		Cluster: aws.String(clusetrName),
	}
	result, err := ecsClient.ListServices(listInput)
	if err != nil {
		return false, err
	}
	for _, v := range result.ServiceArns {
		s, err := ParseArn(v)
		if err != nil {
			return false, err
		}
		if *s.Name == serviceName {
			return true, nil
		}
	}
	return false, nil
}

func ParseArn(arn *string) (*Service, error) {
	splitStr := strings.Split(*arn, "service/")
	if len(splitStr) != 2 {
		return nil, errors.New("illegal arn string")
	}
	return &Service{&splitStr[1]}, nil
}
