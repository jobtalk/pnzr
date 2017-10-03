package setting

import "github.com/aws/aws-sdk-go/service/ecs"

type Setting struct {
	Version float64
	Service        *ecs.CreateServiceInput
	TaskDefinition *ecs.RegisterTaskDefinitionInput
}