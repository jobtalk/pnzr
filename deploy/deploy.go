package deploy

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
	. "github.com/ieee0824/thor/ecs"
)

type Setting struct {
	Service        *ecs.CreateServiceInput
	TaskDefinition *ecs.RegisterTaskDefinitionInput
}

// serviceが存在しない時はサービスを作る
// 存在するときはアップデートする
func Deploy(awsConfig aws.Config, s *Setting) (interface{}, error) {
	if s.TaskDefinition != nil {
		_, err := RegisterTaskDefinition(awsConfig, s.TaskDefinition)
		if err != nil {
			return nil, err
		}
	}
	return UpsertService(awsConfig, s.Service)
}
