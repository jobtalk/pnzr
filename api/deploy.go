package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/jobtalk/eriri/lib"
	"github.com/jobtalk/eriri/lib/setting"
)

// serviceが存在しない時はサービスを作る
// 存在するときはアップデートする
func Deploy(awsConfig *aws.Config, s *setting.Setting) (interface{}, error) {
	ecsClient := lib.NewECS(awsConfig)
	var result = []interface{}{}
	if s.ELB != nil {
		resultMkELB, err := MkELB(awsConfig, s.ELB)
		if err != nil {
			return nil, err
		}
		result = append(result, resultMkELB)
		for _, v := range result {
			if es, ok := v.([]interface{}); ok {
				for _, e := range es {
					if tgOutput, ok := e.(*elbv2.CreateTargetGroupOutput); ok {
						for _, tg := range tgOutput.TargetGroups {
							lb := s.ECS.Service.LoadBalancers[0]
							lb.TargetGroupArn = tg.TargetGroupArn
							s.ECS.Service.LoadBalancers[0] = lb
						}
					}
				}
			}
		}
	}
	if s.ECS != nil && s.ECS.TaskDefinition != nil {
		resultTaskDefinition, err := ecsClient.RegisterTaskDefinition(s.ECS.TaskDefinition)
		if err != nil {
			return nil, err
		}
		result = append(result, resultTaskDefinition)
	}
	resultUpsert, err := lib.UpsertService(awsConfig, s.ECS.Service)
	if err != nil {
		return nil, err
	}

	return append(result, resultUpsert), nil
}
