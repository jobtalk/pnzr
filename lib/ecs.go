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

func NewECS(sess *session.Session) *ECS {
	return &ECS{
		svc: ecs.New(sess),
	}
}

func (e *ECS) RegisterTaskDefinition(registerTaskDefinitionInput *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
	return e.svc.RegisterTaskDefinition(registerTaskDefinitionInput)
}

func (e *ECS) UpsertService(createServiceInput *ecs.CreateServiceInput) (interface{}, error) {
	if createServiceInput.Cluster == nil {
		createServiceInput.Cluster = aws.String("default")
	}
	ok, err := e.serviceExists(*createServiceInput.Cluster, *createServiceInput.ServiceName)
	if err != nil {
		return nil, err
	}
	if !ok {
		return e.svc.CreateService(createServiceInput)
	}

	updateServiceInput := &ecs.UpdateServiceInput{}
	updateServiceInput.Cluster = createServiceInput.Cluster
	updateServiceInput.DeploymentConfiguration = createServiceInput.DeploymentConfiguration
	updateServiceInput.DesiredCount = createServiceInput.DesiredCount
	updateServiceInput.Service = createServiceInput.ServiceName
	updateServiceInput.TaskDefinition = createServiceInput.TaskDefinition
	return e.svc.UpdateService(updateServiceInput)
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

func (e *ECS) serviceExists(clusetrName string, serviceName string) (bool, error) {
	listInput := &ecs.ListServicesInput{
		Cluster: aws.String(clusetrName),
	}
	result, err := e.ListServices(listInput)
	if err != nil {
		return false, err
	}
	for _, v := range result.ServiceArns {
		s, err := parseArn(v)
		if err != nil {
			return false, err
		}
		if *s.Name == serviceName {
			return true, nil
		}
	}
	return false, nil
}

func parseArn(arn *string) (*Service, error) {
	splitStr := strings.Split(*arn, "service/")
	if len(splitStr) != 2 {
		return nil, errors.New("illegal arn string")
	}
	return &Service{&splitStr[1]}, nil
}
