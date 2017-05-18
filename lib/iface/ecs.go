package iface

import (
	"github.com/aws/aws-sdk-go/service/ecs"
)

type ECSAPI interface {
	RegisterTaskDefinition(*ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error)
	UpsertService(*ecs.CreateServiceInput) (interface{}, error)
}
