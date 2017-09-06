package api

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jobtalk/pnzr/lib"
	"github.com/jobtalk/pnzr/lib/config"
	"github.com/jobtalk/pnzr/lib/iface"
)

type DeployDeps struct {
	ecs iface.ECSAPI
}

// serviceが存在しない時はサービスを作る
// 存在するときはアップデートする

func (d *DeployDeps) Deploy(c *config.IntermediateConfig) (interface{}, error) {
	var result = []interface{}{}
	if c != nil && c.TaskDefinition != nil {
		resultTaskDefinition, err := d.ecs.RegisterTaskDefinition(c.TaskDefinition)
		if err != nil {
			return nil, err
		}
		result = append(result, resultTaskDefinition)
	}
	if c != nil && c.Service != nil {
		resultUpsert, err := d.ecs.UpsertService(c.Service)
		if err != nil {
			return nil, err
		}
		result = append(result, resultUpsert)
	}

	return result, nil
}

func Deploy(sess *session.Session, c *config.IntermediateConfig) (interface{}, error) {
	return (&DeployDeps{ecs: lib.NewECS(sess)}).Deploy(c)
}
