package deploy

import (
	"github.com/aws/aws-sdk-go/aws"
	. "github.com/ieee0824/thor/ecs"
	"github.com/ieee0824/thor/mkelb"
	"github.com/ieee0824/thor/setting"
)

// serviceが存在しない時はサービスを作る
// 存在するときはアップデートする

func Deploy(awsConfig *aws.Config, s *setting.Setting) (interface{}, error) {
	var result = []interface{}{}
	if s.ECS != nil {
		resultMkELB, err := mkelb.MkELB(awsConfig, s.ELB)
		if err != nil {
			return nil, err
		}
		result = append(result, resultMkELB)
	}
	if s.ECS.TaskDefinition != nil {
		resultTaskDefinition, err := RegisterTaskDefinition(awsConfig, s.ECS.TaskDefinition)
		if err != nil {
			return nil, err
		}
		result = append(result, resultTaskDefinition)
	}
	resultUpsert, err := UpsertService(awsConfig, s.ECS.Service)
	if err != nil {
		return nil, err
	}

	return append(result, resultUpsert), nil
}
