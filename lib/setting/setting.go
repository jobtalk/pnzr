package setting

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

type ELB struct {
	LB          *elbv2.CreateLoadBalancerInput
	TargetGroup *elbv2.CreateTargetGroupInput
	Listener    *elbv2.CreateListenerInput
}

type ECS struct {
	Service        *ecs.CreateServiceInput
	TaskDefinition *ecs.RegisterTaskDefinitionInput
}

type Setting struct {
	ELB *ELB
	ECS *ECS
}
