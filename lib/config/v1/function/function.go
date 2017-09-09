package function

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
	"path/filepath"
)

type JSFunction struct {
	parentSettingPath string
}

func New(p string) *JSFunction {
	return &JSFunction{
		parentSettingPath: p,
	}
}

func (f *JSFunction)Require(call otto.FunctionCall) otto.Value {
	dir := filepath.Dir(f.parentSettingPath)
	file := call.Argument(0).String()
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, file))
	if err != nil {
		panic(err)
	}
	_, err = call.Otto.Run(string(data))
	if err != nil {
		panic(err)
	}
	return call.This
}

func (f *JSFunction)LoadJSON(call otto.FunctionCall) otto.Value {
	dir := filepath.Dir(f.parentSettingPath)
	file := call.Argument(0).String()
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", dir, file))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	m := map[string]interface{}{}
	if err := json.Unmarshal(data, &m); err != nil {
		panic(err)
	}

	v, err := call.Otto.ToValue(m)
	if err != nil {
		panic(err)
	}

	return v
}
