package api

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jobtalk/pnzr/lib"
	"github.com/jobtalk/pnzr/lib/iface"
	"github.com/jobtalk/pnzr/lib/setting"
)

type DeployDeps struct {
	ecs iface.ECSAPI
}

// serviceが存在しない時はサービスを作る
// 存在するときはアップデートする
func (d *DeployDeps) Deploy(s *setting.Setting) (interface{}, error) {
	var result = []interface{}{}
	if s != nil && s.TaskDefinition != nil {
		resultTaskDefinition, err := d.ecs.RegisterTaskDefinition(s.TaskDefinition)
		if err != nil {
			return nil, err
		}
		result = append(result, resultTaskDefinition)
	}
	if s != nil && s.Service != nil {
		resultUpsert, err := d.ecs.UpsertService(s.Service)
		if err != nil {
			return nil, err
		}
		result = append(result, resultUpsert)
	}

	return result, nil
}

func Deploy(sess *session.Session, s *setting.Setting) (interface{}, error) {
	return (&DeployDeps{ecs: lib.NewECS(sess)}).Deploy(s)
}
