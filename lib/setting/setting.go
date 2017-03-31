package setting

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type ELB struct {
	*elbv2.CreateLoadBalancerInput
	*elbv2.CreateTargetGroupInput
	*elbv2.CreateListenerInput
}

type ECS struct {
	Service        *ecs.CreateServiceInput
	TaskDefinition *ecs.RegisterTaskDefinitionInput
}

type Setting struct {
	ELB *ELB
	ECS *ECS
}
