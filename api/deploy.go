package api

import (
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
	if s.ECS != nil && s.ECS.TaskDefinition != nil {
		resultTaskDefinition, err := d.ecs.RegisterTaskDefinition(s.ECS.TaskDefinition)
		if err != nil {
			return nil, err
		}
		result = append(result, resultTaskDefinition)
	}
	if s.ECS != nil && s.ECS.Service != nil {
		resultUpsert, err := d.ecs.UpsertService(s.ECS.Service)
		if err != nil {
			return nil, err
		}
		result = append(result, resultUpsert)
	}

	return result, nil
}

func Deploy(s *setting.Setting) (interface{}, error) {
	return (&DeployDeps{ecs: lib.NewECS()}).Deploy(s)
}
