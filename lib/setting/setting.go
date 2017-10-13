package setting

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/service/ecs"
	"io/ioutil"
)

type Setting struct {
	Version        float64
	Service        *ecs.CreateServiceInput
	TaskDefinition *ecs.RegisterTaskDefinitionInput
}

func IsV1Setting(path string) bool {
	bin, err := ioutil.ReadFile(path)
	if err != nil {
		return false
	}
	m := map[string]interface{}{}
	if err := json.Unmarshal(bin, &m); err != nil {
		return false
	}

	v, ok := m["Version"].(float64)
	if !ok {
		return false
	}

	return v == 1.0
}
