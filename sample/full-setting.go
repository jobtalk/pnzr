package sample

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/ieee0824/thor/setting"
)

func sample() string {
	s := &setting.Setting{
		ELB: &setting.ELB{
			&elbv2.CreateLoadBalancerInput{
				IpAddressType: aws.String("'ipv4' or 'dualstack'"),
				Name:          aws.String("load balancer name"),
				Scheme:        aws.String("https://github.com/aws/aws-sdk-go/blob/master/service/elbv2/api.go#L2798"),
				SecurityGroups: []*string{
					aws.String("sg-xxxxxx"),
					aws.String("sg-yyyyyy"),
				},
				Subnets: []*string{
					aws.String("subnet-xxxxxx"),
					aws.String("subnet-yyyyyy"),
				},
				Tags: []*elbv2.Tag{
					&elbv2.Tag{
						Key:   aws.String("key_0"),
						Value: aws.String("foo"),
					},
					&elbv2.Tag{
						Key:   aws.String("key_1"),
						Value: aws.String("bar"),
					},
				},
			},
			&elbv2.CreateTargetGroupInput{
				HealthCheckIntervalSeconds: aws.Int64(5),
				HealthCheckPath:            aws.String("/"),
				HealthCheckPort:            aws.String("80"),
				HealthCheckProtocol:        aws.String("HTTP"),
				HealthCheckTimeoutSeconds:  aws.Int64(5),
				HealthyThresholdCount:      aws.Int64(5),
				Matcher:                    &elbv2.Matcher{},
				Name:                       aws.String("target group name"),
				Port:                       aws.Int64(80),
				Protocol:                   aws.String("HTTP"),
				UnhealthyThresholdCount:    aws.Int64(2),
				VpcId: aws.String("vpc-xxxxxx"),
			},
			&elbv2.CreateListenerInput{
				Certificates: []*elbv2.Certificate{
					&elbv2.Certificate{
						CertificateArn: aws.String("certifacate arn"),
					},
				},
				DefaultActions: []*elbv2.Action{
					&elbv2.Action{},
				},
				LoadBalancerArn: nil,
			},
		},
		ECS: &setting.ECS{
			&ecs.CreateServiceInput{
				DeploymentConfiguration: &ecs.DeploymentConfiguration{},
				LoadBalancers: []*ecs.LoadBalancer{
					&ecs.LoadBalancer{},
				},
				PlacementConstraints: []*ecs.PlacementConstraint{
					&ecs.PlacementConstraint{},
				},
				PlacementStrategy: []*ecs.PlacementStrategy{
					&ecs.PlacementStrategy{},
				},
			},
			&ecs.RegisterTaskDefinitionInput{
				ContainerDefinitions: []*ecs.ContainerDefinition{
					&ecs.ContainerDefinition{},
				},
				PlacementConstraints: []*ecs.TaskDefinitionPlacementConstraint{
					&ecs.TaskDefinitionPlacementConstraint{},
				},
				Volumes: []*ecs.Volume{
					&ecs.Volume{},
				},
			},
		},
	}

	bin, _ := json.MarshalIndent(s, "", "    ")
	return string(bin)
}
