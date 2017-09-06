package config

import "github.com/aws/aws-sdk-go/service/ecs"

type IntermediateConfig struct {
	Version        float64
	Service        *ecs.CreateServiceInput
	TaskDefinition *ecs.RegisterTaskDefinitionInput
}
