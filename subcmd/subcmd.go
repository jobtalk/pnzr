package subcmd

import (
	"errors"
	"fmt"
	"strings"
)

// --hoge=hugaみたいなやつ
func getFullNameParam(args []string, key string) ([]*string, error) {
	var result = []*string{}
	for _, v := range args {
		if strings.Contains(v, key) {
			splitStr := strings.Split(v, "=")
			if len(splitStr) == 1 {
				param := "true"
				result = append(result, &param)
			} else if len(splitStr) != 2 {
				return nil, errors.New(fmt.Sprintf("%s is illegal parameter", key))
			} else if splitStr[0] == key {
				result = append(result, &splitStr[1])
			}
		}
	}
	return result, nil
}

// -f hogeみたいなやつ
func getValFromArgs(args []string, key string) ([]*string, error) {
	var result = []*string{}
	for i, v := range args {
		if v == key {
			// vが一番最後じゃないとき
			if i+1 != len(args) {
				result = append(result, &args[i+1])
			} else {
				return nil, errors.New(fmt.Sprintf("%s is illegal parameter", key))
			}
		}
	}
	return result, nil
}
