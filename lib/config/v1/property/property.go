package property

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/robertkrimen/otto"
	"encoding/json"
)

type DeployConfigration struct {
	Service        *ecs.CreateServiceInput
	TaskDefinition *ecs.RegisterTaskDefinitionInput
}

type Property struct {
	Version float64
	Service  map[string]interface{}
	TaskDefinition map[string]interface{}
}

func (p *Property) convertService() (*ecs.CreateServiceInput, error) {
	var service = ecs.CreateServiceInput{}
	body, ok := p.Service["body"]
	if !ok {
		return nil, nil
	}
	v, ok := body.(otto.Value)
	if !ok {

		return nil, nil
	}
	i, err := v.Export()
	if err != nil {
		return nil, err
	}
	bin, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bin, &service); err != nil {
		return nil, err
	}
	return &service, nil
}

func (p *Property) convertTaskDefinition() (*ecs.RegisterTaskDefinitionInput, error) {
	var taskDefinition = ecs.RegisterTaskDefinitionInput{}

	body, ok := p.TaskDefinition["body"]
	if !ok {
		return nil, nil
	}
	v, ok := body.(otto.Value)
	if !ok {
		return nil, nil
	}
	i, err := v.Export()
	if err != nil {
		return nil, err
	}
	bin, err := json.Marshal(i)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bin, &taskDefinition); err != nil {
		return nil, err
	}
	return &taskDefinition, nil
}

func (p *Property) ConvertToConfig() (*DeployConfigration, error) {
	ret := DeployConfigration{}

	service, err := p.convertService()
	if err != nil {
		return nil, err
	}
	taskDefinition, err := p.convertTaskDefinition()
	if err != nil {
		return nil, err
	}

	ret.Service = service
	ret.TaskDefinition = taskDefinition

	return &ret, nil
}

