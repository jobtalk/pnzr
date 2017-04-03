package api

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jobtalk/thor/lib"
	"github.com/jobtalk/thor/lib/setting"
)

// serviceが存在しない時はサービスを作る
// 存在するときはアップデートする
func Deploy(awsConfig *aws.Config, s *setting.Setting) (interface{}, error) {
	var result = []interface{}{}
	if s.ECS != nil {
		resultMkELB, err := MkELB(awsConfig, s.ELB)
		if err != nil {
			return nil, err
		}
		result = append(result, resultMkELB)
	}
	if s.ECS.TaskDefinition != nil {
		resultTaskDefinition, err := lib.RegisterTaskDefinition(awsConfig, s.ECS.TaskDefinition)
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
