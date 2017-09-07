package function

import (
	"encoding/json"
	"fmt"
	"github.com/robertkrimen/otto"
	"io/ioutil"
)

func Require(call otto.FunctionCall) otto.Value {
	file := call.Argument(0).String()
	data, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	_, err = call.Otto.Run(string(data))
	if err != nil {
		panic(err)
	}
	return call.This
}

func LoadJSON(call otto.FunctionCall) otto.Value {
	file := call.Argument(0).String()
	data, err := ioutil.ReadFile(file)
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
