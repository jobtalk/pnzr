package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/jobtalk/thor/lib"
	"github.com/jobtalk/thor/lib/setting"
)

func createTargetGroup(awsConfig *aws.Config, s *elbv2.CreateTargetGroupInput) (*elbv2.CreateTargetGroupOutput, error) {
	resultTargetGroup, err := lib.CreateTargetGroup(awsConfig, s)
	if err != nil {
		return nil, err
	}
	return resultTargetGroup, nil
}

func createLoadBalancer(awsConfig *aws.Config, s *elbv2.CreateLoadBalancerInput) (*elbv2.CreateLoadBalancerOutput, error) {
	r, err := lib.CreateLoadBalancer(awsConfig, s)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func MkELB(awsConfig *aws.Config, s *setting.ELB) (interface{}, error) {
	var result = []interface{}{}
	var resultTargetGroup *elbv2.CreateTargetGroupOutput
	var resultLoadBalancer *elbv2.CreateLoadBalancerOutput

	if s.TargetGroup != nil {
		r, err := createTargetGroup(awsConfig, s.TargetGroup)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
		resultTargetGroup = r
	}

	if s.LB != nil {
		r, err := createLoadBalancer(awsConfig, s.LB)
		if err != nil {
			return nil, err
		}
		result = append(result, r)
		resultLoadBalancer = r
	}

	defaultAction := &elbv2.Action{
		TargetGroupArn: resultTargetGroup.TargetGroups[0].TargetGroupArn,
		Type:           aws.String("forward"),
	}

	s.Listener.DefaultActions = []*elbv2.Action{
		defaultAction,
	}

	s.Listener.LoadBalancerArn = resultLoadBalancer.LoadBalancers[0].LoadBalancerArn

	resultLister, err := lib.CreateListener(awsConfig, s.Listener)
	if err != nil {
		return nil, err
	}
	result = append(result, resultLister)
	return result, nil
}
