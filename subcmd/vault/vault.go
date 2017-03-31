package vault

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/jobtalk/thor/api"
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

type vaultParam struct {
	Pass *string
	Path *string
}

func (v *vaultParam) validate() error {
	if v.Pass == nil {
		return errors.New("pass is empty")
	} else if v.Path == nil {
		return errors.New("path is empty")
	} else if *v.Pass == "" {
		return errors.New("pass is empty")
	} else if *v.Path == "" {
		return errors.New("path is empty")
	}
	return nil
}

func parseVaultArgs(args []string) (*vaultParam, error) {
	var result = &vaultParam{}
	passPram, err := getValFromArgs(args, "-p")
	if err != nil {
		return nil, err
	}
	if 2 <= len(passPram) {
		return nil, errors.New("'-p' parameter can not be specified more than once.")
	} else if 0 == len(passPram) {
		return nil, errors.New("'-p' parameter is empty")
	}
	result.Pass = passPram[0]

	pathParam, err := getValFromArgs(args, "-f")
	if err != nil {
		return nil, err
	}
	if 2 <= len(pathParam) {
		return nil, errors.New("'-f' parameter can not be specified more than once.")
	} else if 0 == len(pathParam) {
		return nil, errors.New("'-f' parameter is empty")
	}
	result.Path = pathParam[0]
	return result, nil
}

type Vault struct{}

func (c *Vault) Help() string {
	help := ""
	help += "usage: vault [options ...]\n"
	help += "options:\n"
	help += "    -f vault target json\n"
	help += "\n"
	help += "    -p vault pass\n"

	return help
}

func (c *Vault) Run(args []string) int {
	param, err := parseVaultArgs(args)
	if err != nil {
		log.Fatalln(err)
	}
	if err := param.validate(); err != nil {
		log.Fatalln(err)
	}
	bin, err := ioutil.ReadFile(*param.Path)
	if err != nil {
		log.Fatalln(err)
	}
	vaulter := api.New(bin)
	if err := vaulter.Encrypt(*param.Pass); err != nil {
		log.Fatalln(err)
	}
	vaultedJSON, err := json.Marshal(vaulter)
	if err != nil {
		log.Fatalln(err)
	}

	if err := ioutil.WriteFile(*param.Path, vaultedJSON, 0644); err != nil {
		log.Fatalln(err)
	}

	return 0
}

func (c *Vault) Synopsis() string {
	synopsis := ""
	synopsis += "usage: thor vault [options ...]\n"
	synopsis += "options:\n"
	synopsis += "    -f vault target json\n"
	synopsis += "\n"
	synopsis += "    -p vault password\n"
	synopsis += "===================================================\n"

	return synopsis
}
