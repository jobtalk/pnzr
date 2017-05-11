package lib

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/aws/aws-sdk-go/service/elbv2/elbv2iface"
)

type ELB struct {
	svc elbv2iface.ELBV2API
}

func NewELB(awsConfig *aws.Config) *ELB {
	return &ELB{
		svc: elbv2.New(session.New(), awsConfig),
	}
}

func (e *ELB) CreateLoadBalancer(createLoadBalancerInput *elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
	return e.svc.CreateLoadBalancer(createLoadBalancerInput)
}

func (e *ELB) CreateTargetGroup(params *elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
	ret, err := e.svc.CreateTargetGroup(params)
	if err != nil {
		return nil, err
	}
	targetGroupARN = ret.TargetGroups[0].TargetGroupArn
	return ret, nil
}

func (e *ELB) CreateListener(params *elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
	return e.svc.CreateListener(params)
}

func (e *ELB) DescribeLoadBalancers(params *elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
	return e.svc.DescribeLoadBalancers(params)
}
