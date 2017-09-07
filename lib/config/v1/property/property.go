package property

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"encoding/json"
)

type Property struct {
	Version        float64
	Config interface{}
}

func (p *Property)ConvertToConfig() (*DeployConfigration, error) {
	bin, err := json.Marshal(p.Config)
	if err != nil {
		return nil, err
	}
	ret := &DeployConfigration{}
	if err := json.Unmarshal(bin, ret); err != nil {
		return nil, err
	}

	return ret, nil
}
type DeployConfigration struct {
	Service        *ecs.CreateServiceInput
	TaskDefinition *ecs.RegisterTaskDefinitionInput
}
