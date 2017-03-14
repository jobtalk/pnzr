package subcmd

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/ieee0824/thor/vault"
)

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
	vaulter := vault.New(bin)
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
