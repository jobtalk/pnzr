package v1_config

import (
	"github.com/jobtalk/pnzr/lib/config"
	"github.com/jobtalk/pnzr/lib/config/v1/function"
	"github.com/jobtalk/pnzr/lib/config/v1/property"
	"github.com/robertkrimen/otto"
)

const (
	VERSION = 1.0
)

func CheckSupportVersion(confPath *string) bool {
	vm := otto.New()
	prop := property.Property{}

	script, err := vm.Compile(*confPath, nil)
	if err != nil {
		return false
	}
	vm.Set("config", &prop)
	if _, err := vm.Run(script); err != nil {
		return false
	}
	return prop.Version == VERSION
}

type ConfigLoader struct {
}

func (c *ConfigLoader) Load(confPath *string) (*config.IntermediateConfig, error) {
	prop := property.Property{}
	vm := otto.New()

	script, err := vm.Compile(*confPath, nil)
	if err != nil {
		return nil, err
	}

	vm.Set("require", function.Require)
	vm.Set("loadJSON", function.LoadJSON)
	vm.Set("config", &prop)

	if _, err := vm.Run(script); err != nil {
		return nil, err
	}

	deployConfig, err := prop.ConvertToConfig()
	if err != nil {
		return nil, err
	}

	return &config.IntermediateConfig{Service: deployConfig.Service, TaskDefinition: deployConfig.TaskDefinition, Version: VERSION}, nil
}
