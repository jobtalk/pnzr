package subcmd

import "errors"

type vaultParam struct {
	Pass *string
	Path *string
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
	return ""
}

func (c *Vault) Run(args []string) int {
	return 0
}

func (c *Deploy) Synopsis() string {
	return ""
}
