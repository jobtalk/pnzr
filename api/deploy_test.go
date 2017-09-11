package api

import (
	"testing"

	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/jobtalk/pnzr/lib/config"
	"github.com/jobtalk/pnzr/lib/iface"
)

func TestDeploy(t *testing.T) {
	family := "taskdef-a"
	taskDefinition := &ecs.RegisterTaskDefinitionInput{Family: &family}
	serviceName := "service-a"
	service := &ecs.CreateServiceInput{ServiceName: &serviceName}

	// Settings
	onlyService := config.IntermediateConfig{
		Service: service,
	}
	onlyTaskDefinition := config.IntermediateConfig{
		TaskDefinition: taskDefinition,
	}
	both := config.IntermediateConfig{
		Service:        service,
		TaskDefinition: taskDefinition,
	}

	// UpsertService のみが呼ばれる
	{
		deploy, fnArgs := mockDeploy()
		deploy.Deploy(&onlyService)
		if fnArgs.RegisterTaskDefinitionInput != nil {
			t.Fatalf("RegisterTaskDefinition should not be called")
		}
		if *fnArgs.UpsertServiceInput.ServiceName != serviceName {
			t.Log(*fnArgs.UpsertServiceInput.ServiceName)
			t.Fatalf("UpsertService should be called with %s", serviceName)
		}
	}

	// RegisterTaskDefinition のみが呼ばれる
	{
		deploy, fnArgs := mockDeploy()
		deploy.Deploy(&onlyTaskDefinition)
		if *fnArgs.RegisterTaskDefinitionInput.Family != family {
			t.Log(*fnArgs.RegisterTaskDefinitionInput.Family)
			t.Fatalf("RegisterTaskDefinition should be called with %s", family)
		}
		if fnArgs.UpsertServiceInput != nil {
			t.Fatalf("UpsertService should not be called")
		}
	}

	// UpsertService, RegisterTaskDefinition 両方が呼ばれる
	{
		deploy, fnArgs := mockDeploy()
		deploy.Deploy(&both)
		if *fnArgs.RegisterTaskDefinitionInput.Family != family {
			t.Log(*fnArgs.RegisterTaskDefinitionInput.Family)
			t.Fatalf("RegisterTaskDefinition should be called with %s", family)
		}
		if *fnArgs.UpsertServiceInput.ServiceName != serviceName {
			t.Log(*fnArgs.UpsertServiceInput.ServiceName)
			t.Fatalf("UpsertService should be called with %s", serviceName)
		}
	}
}

type mockedFnArgs struct {
	RegisterTaskDefinitionInput *ecs.RegisterTaskDefinitionInput
	UpsertServiceInput          *ecs.CreateServiceInput
}

func mockDeploy() (*DeployDeps, *mockedFnArgs) {
	fnArgs := mockedFnArgs{RegisterTaskDefinitionInput: nil, UpsertServiceInput: nil}
	ecs := mockedECS{
		registerTaskDefinition: func(in *ecs.RegisterTaskDefinitionInput) (*ecs.RegisterTaskDefinitionOutput, error) {
			fnArgs = mockedFnArgs{RegisterTaskDefinitionInput: in, UpsertServiceInput: fnArgs.UpsertServiceInput}
			return nil, nil
		},
		upsertService: func(in *ecs.CreateServiceInput) (interface{}, error) {
			fnArgs = mockedFnArgs{RegisterTaskDefinitionInput: fnArgs.RegisterTaskDefinitionInput, UpsertServiceInput: in}
			return nil, nil
		},
	}
	return &DeployDeps{ecs: ecs}, &fnArgs
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
