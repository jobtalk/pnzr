package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"strings"
)

func parse(v string) (key string, val interface{}, err error) {
	s := strings.SplitN(v, "=",1)
	if len(s) != 2 {
		return "", nil, fmt.Errorf("parse error: %v", v)
	}
	key = s[0]
	if err = json.Unmarshal([]byte(s[1]), &val); err != nil {
		return "", nil, err
	}
	return
}

func main() {
	m := map[string]interface{}{}
	if len(os.Args) != 3 {
		panic(fmt.Errorf("illegal args: %v", os.Args))
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	f.Close()
	if err := json.Unmarshal(body, &m); err != nil {
		panic(err)
	}

	key, val, err := parse(os.Args[2])
	if err != nil {
		panic(err)
	}
	m[key] = val
	result, err := json.MarshalIndent(m, "", "    ")
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(os.Args[1], result, 0644); err != nil {
		panic(err)
	}
}
