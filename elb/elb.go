package elb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

func CreateLoadBalancer(awsConfig *aws.Config, createLoadBalancerInput *elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
	svc := elbv2.New(session.New(), awsConfig)

	return svc.CreateLoadBalancer(createLoadBalancerInput)
}

func CreateTargetGroup(awsConfig *aws.Config, params *elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
	svc := elbv2.New(session.New(), awsConfig)
	return svc.CreateTargetGroup(params)
}

func CreateListener(awsConfig *aws.Config, params *elbv2.CreateListenerInput) (*elbv2.CreateListenerOutput, error) {
	svc := elbv2.New(session.New(), awsConfig)
	return svc.CreateListener(params)
}

func DescribeLoadBalancers(awsConfig *aws.Config, params *elbv2.DescribeLoadBalancersInput) (*elbv2.DescribeLoadBalancersOutput, error) {
	svc := elbv2.New(session.New(), awsConfig)
	return svc.DescribeLoadBalancers(params)
}
