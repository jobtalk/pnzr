package mkelb

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/ieee0824/thor/elb"
	"github.com/ieee0824/thor/setting"
)

func MkELB(awsConfig *aws.Config, s *setting.ELB) (interface{}, error) {
	var result = []interface{}{}
	resultTargetGroup, err := elb.CreateTargetGroup(awsConfig, s.CreateTargetGroupInput)
	if err != nil {
		return nil, err
	}
	result = append(result, resultTargetGroup)

	resultLoadBalancer, err := elb.CreateLoadBalancer(awsConfig, s.CreateLoadBalancerInput)
	if err != nil {
		return nil, err
	}
	result = append(result, resultLoadBalancer)

	defaultAction := &elbv2.Action{
		TargetGroupArn: resultTargetGroup.TargetGroups[0].TargetGroupArn,
		Type:           aws.String("forward"),
	}
	s.CreateListenerInput.DefaultActions = []*elbv2.Action{
		defaultAction,
	}
	s.CreateListenerInput.LoadBalancerArn = resultLoadBalancer.LoadBalancers[0].LoadBalancerArn
	resultLister, err := elb.CreateListener(awsConfig, s.CreateListenerInput)
	if err != nil {
		return nil, err
	}
	result = append(result, resultLister)
	return result, nil
}
