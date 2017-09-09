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
	var service = map[string]interface{}{}
	var taskDefinition = map[string]interface{}{}

	vm := otto.New()
	prop := property.Property{}

	script, err := vm.Compile(*confPath, nil)
	if err != nil {
		return false
	}
	f := function.New(*confPath)
	vm.Set("require", f.Require)
	vm.Set("loadJSON", f.LoadJSON)
	vm.Set("config", &prop)
	vm.Set("service", service)
	vm.Set("taskDefinition", taskDefinition)
	if _, err := vm.Run(script); err != nil {
		return false
	}

	return prop.Version == VERSION
}

type ConfigLoader struct {
}

func (c *ConfigLoader) Load(confPath *string) (*config.IntermediateConfig, error) {
	var service = map[string]interface{}{}
	var taskDefinition = map[string]interface{}{}

	prop := property.Property{}
	vm := otto.New()

	script, err := vm.Compile(*confPath, nil)
	if err != nil {
		return nil, err
	}

	f := function.New(*confPath)
	vm.Set("require", f.Require)
	vm.Set("loadJSON", f.LoadJSON)
	vm.Set("config", &prop)
	vm.Set("service", service)
	vm.Set("taskDefinition", taskDefinition)

	if _, err := vm.Run(script); err != nil {
		return nil, err
	}

	prop.Service = service
	prop.TaskDefinition = taskDefinition

	deployConfig, err := prop.ConvertToConfig()
	if err != nil {
		return nil, err
	}

	return &config.IntermediateConfig{Service: deployConfig.Service, TaskDefinition: deployConfig.TaskDefinition, Version: VERSION}, nil
}
