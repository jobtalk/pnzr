package property

import "github.com/aws/aws-sdk-go/service/ecs"

type Property struct {
	Version        float64
	Service        *ecs.CreateServiceInput
	TaskDefinition *ecs.RegisterTaskDefinitionInput
}
