package lib

import (
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

type Service struct {
	Name *string
}

func RegisterTaskDefinition(awsConfig *aws.Config, registerTaskDefinitionInput *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
	svc := ecs.New(session.New(), awsConfig)
	return svc.RegisterTaskDefinition(registerTaskDefinitionInput)
}

func ListServices(awsConfig *aws.Config, params *ecs.ListServicesInput) (*ecs.ListServicesOutput, error) {
	var (
		ret     = &ecs.ListServicesOutput{}
		pageNum int
		svc     = ecs.New(session.New(), awsConfig)
	)

	err := svc.ListServicesPages(params, func(page *ecs.ListServicesOutput, lastPage bool) bool {
		pageNum++
		ret.ServiceArns = append(ret.ServiceArns, page.ServiceArns...)
		return pageNum <= 1000
	})
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func CreateService(awsConfig *aws.Config, createServiceInput *ecs.CreateServiceInput) (*ecs.CreateServiceOutput, error) {
	svc := ecs.New(session.New(), awsConfig)
	return svc.CreateService(createServiceInput)
}

func UpdateService(awsConfig *aws.Config, updateServiceInput *ecs.UpdateServiceInput) (*ecs.UpdateServiceOutput, error) {
	svc := ecs.New(session.New(), awsConfig)
	return svc.UpdateService(updateServiceInput)
}

func UpsertService(awsConfig *aws.Config, createServiceInput *ecs.CreateServiceInput) (interface{}, error) {
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
		return CreateService(awsConfig, createServiceInput)
	}

	updateServiceInput := &ecs.UpdateServiceInput{}
	updateServiceInput.Cluster = createServiceInput.Cluster
	updateServiceInput.DeploymentConfiguration = createServiceInput.DeploymentConfiguration
	updateServiceInput.DesiredCount = createServiceInput.DesiredCount
	updateServiceInput.Service = createServiceInput.ServiceName
	updateServiceInput.TaskDefinition = createServiceInput.TaskDefinition
	return UpdateService(awsConfig, updateServiceInput)
}

func IsExistService(awsConfig *aws.Config, clusetrName string, serviceName string) (bool, error) {
	listInput := &ecs.ListServicesInput{
		Cluster: aws.String(clusetrName),
	}
	result, err := ListServices(awsConfig, listInput)
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
