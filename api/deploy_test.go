package api

import (
	_ "fmt"
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/jobtalk/pnzr/lib/iface"
	"github.com/jobtalk/pnzr/lib/setting"
)

func TestDeploy(t *testing.T) {
	taskDefinition := &ecs.RegisterTaskDefinitionInput{}
	service := &ecs.CreateServiceInput{}

	// Settings
	onlyService := setting.Setting{
		ECS: &setting.ECS{Service: service},
	}
	onlyTaskDefinition := setting.Setting{
		ECS: &setting.ECS{TaskDefinition: taskDefinition},
	}
	both := setting.Setting{
		ECS: &setting.ECS{
			Service:        service,
			TaskDefinition: taskDefinition,
		},
	}

	// onlyService
	{
		deploy, isRegisterTaskDefinitionCalled, isUpsertServiceCalled := mockDeploy()
		deploy.Deploy(&onlyService)
		if *isRegisterTaskDefinitionCalled {
			t.Fatalf("RegisterTaskDefinition should not be called")
		}
		if !*isUpsertServiceCalled {
			t.Fatalf("UpsertService should be called")
		}
	}

	// onlyTaskDefinition
	{
		deploy, isRegisterTaskDefinitionCalled, isUpsertServiceCalled := mockDeploy()
		deploy.Deploy(&onlyTaskDefinition)
		if !*isRegisterTaskDefinitionCalled {
			t.Fatalf("RegisterTaskDefinition should be called")
		}
		if *isUpsertServiceCalled {
			t.Fatalf("UpsertService should not be called")
		}
	}

	// both
	{
		deploy, isRegisterTaskDefinitionCalled, isUpsertServiceCalled := mockDeploy()
		deploy.Deploy(&both)
		if !*isRegisterTaskDefinitionCalled {
			t.Fatalf("RegisterTaskDefinition should be called")
		}
		if !*isUpsertServiceCalled {
			t.Fatalf("UpsertService should not be called")
		}
	}
}

type mockedECS struct {
	iface.ECSAPI
	registerTaskDefinition func(*ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error)
	upsertService          func(*ecs.CreateServiceInput) (interface{}, error)
}

func (m mockedECS) RegisterTaskDefinition(in *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
	m.registerTaskDefinition(in)
	return nil, nil
}

func (m mockedECS) UpsertService(in *ecs.CreateServiceInput) (interface{}, error) {
	m.upsertService(in)
	return nil, nil
}

func mockDeploy() (*DeployDeps, *bool, *bool) {
	isRegisterTaskDefinitionCalled := false
	isUpsertServiceCalled := false
	ecs := mockedECS{
		registerTaskDefinition: func(in *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
			isRegisterTaskDefinitionCalled = true
			return nil, nil
		},
		upsertService: func(in *ecs.CreateServiceInput) (interface{}, error) {
			isUpsertServiceCalled = true
			return nil, nil
		},
	}
	return &DeployDeps{ecs: ecs}, &isRegisterTaskDefinitionCalled, &isUpsertServiceCalled
}
