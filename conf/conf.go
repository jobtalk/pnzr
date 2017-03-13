package conf

import (
	"encoding/json"
	"fmt"
	"strings"
)

func Embedde(base, val string) (string, error) {
	// 埋め込み用の値をjsonからデコードする
	var v = map[string]interface{}{}
	if err := json.Unmarshal([]byte(val), &v); err != nil {
		return "", err
	}

	// 埋め込み用の値をkeyとvalに分ける
	for k, v := range v {
		// 値を再びjsonに戻す
		valJSON, err := json.Marshal(v)
		if err != nil {
			return "", err
		}

		base = strings.Replace(base, fmt.Sprintf("$%s", k), string(valJSON), -1)
	}
	return base, nil
}
