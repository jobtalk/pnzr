package subcmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/jobtalk/thor/conf"
	"github.com/jobtalk/thor/deploy"
	"github.com/jobtalk/thor/setting"
	"github.com/jobtalk/thor/vault"
)

type deployConfigure struct {
	*setting.Setting
}

type Deploy struct{}

type deployParam struct {
	File    *string
	Profile *string
	Vault   *string
}

func (p deployParam) String() string {
	bin, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(bin)
}

func parseDeployArgs(args []string) (*deployParam, error) {
	var result = &deployParam{}
	/*
		設定ファイルの場所を定義したargsを読む
	*/
	fileParam, err := getValFromArgs(args, "-f")
	if err != nil {
		return nil, err
	} else if len(fileParam) >= 2 {
		return nil, errors.New("'-f' parameter can not be specified more than once.")
	}
	if len(fileParam) == 1 {
		result.File = fileParam[0]
	} else if len(fileParam) == 0 {
		fileName := "thor.json"
		result.File = &fileName
	}
	var vaultPass string
	/* vaultのpass */
	if bin, err := ioutil.ReadFile(".vault"); err == nil {
		vaultPass = string(bin)
	}
	/*
		--vault-password-file
	*/
	if vaultFileParam, err := getFullNameParam(args, "--vault-password-file"); err == nil {
		if len(vaultFileParam) == 1 {
			if bin, err := ioutil.ReadFile(*vaultFileParam[0]); err == nil {
				vaultPass = string(bin)
			}
		} else {
			return nil, errors.New("--vault-password-file param is invalid")
		}
	}
	/*
		--ask-vault-pass
	*/
	if vaultPassParam, err := getFullNameParam(args, "--ask-vault-pass"); err == nil {
		if len(vaultPassParam) == 1 {
			vaultPass = *vaultPassParam[0]
		} else {
			return nil, errors.New("--ask-vault-pass param is invalid")
		}
	}
	if vaultPass != "" {
		result.Vault = &vaultPass
	}

	/*
		awsのprofileの定義関係
	*/
	profileParam, err := getFullNameParam(args, "--profile")
	if err != nil {
		return nil, err
	} else if len(profileParam) >= 2 {
		return nil, errors.New("'--profile' parameter can not be specified more than once.")
	}
	if len(profileParam) == 1 {
		result.Profile = profileParam[0]
	}

	return result, nil
}

func (c *Deploy) Help() string {
	help := ""
	help += "usage: deploy [options ...]\n"
	help += "options:\n"
	help += "    -f thor_setting.json\n"
	help += "\n"
	help += "    --profile=${aws profile name}\n"
	help += "        --profile option is arbitrary parameter.\n"
	help += "    --vault-password-file=${vault pass file}\n"
	help += "\n"
	help += "    --ask-vault-pass=${vault pass string}\n"

	return help
}

func readExternalVariablesFromFile(path string) ([][]byte, error) {
	var result = [][]byte{}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	infos, err := ioutil.ReadDir(abs)
	if err != nil {
		return nil, err
	}

	for _, info := range infos {
		if !info.IsDir() {
			bin, err := ioutil.ReadFile(abs + "/" + info.Name())
			if err != nil {
				return nil, err
			}
			result = append(result, bin)
		}
	}

	return result, nil
}

// 分離されたファイルを読む
func readExternalVariables(pathes ...string) ([][]byte, error) {
	var result = [][]byte{}
	if len(pathes) == 0 {
		return readExternalVariablesFromFile("./externals")
	}
	for _, path := range pathes {
		r, err := readExternalVariablesFromFile(path)
		if err != nil {
			return nil, err
		}
		result = append(result, r...)
	}
	return result, nil
}

func readConf(params *deployParam) (*deployConfigure, error) {
	var config = &deployConfigure{}
	deployConfigureJSON, err := ioutil.ReadFile(*params.File)
	if err != nil {
		return nil, err
	}
	externals, err := readExternalVariables()
	if err != nil {
		return nil, err
	}
	if len(externals) != 0 {
		base := string(deployConfigureJSON)
		for _, external := range externals {
			if vault.IsSecret(external) {
				if params.Vault == nil {
					return nil, errors.New("vault pass is empty")
				}
				plain, err := vault.Decrypt(external, *params.Vault)
				if err != nil {
					return nil, err
				}
				base, err = conf.Embedde(base, string(plain))
				if err != nil {
					return nil, err
				}
			} else {
				var err error
				base, err = conf.Embedde(base, string(external))
				if err != nil {
					return nil, err
				}
			}
		}
		deployConfigureJSON = []byte(base)
	}

	if err := json.Unmarshal(deployConfigureJSON, config); err != nil {
		return nil, err
	}
	return config, err
}

func (c *Deploy) Run(args []string) int {
	params, err := parseDeployArgs(args)
	if err != nil {
		log.Fatalln(err)
	}
	var cred *credentials.Credentials
	if params.Profile != nil {
		cred = credentials.NewSharedCredentials("", "default")
	}
	awsConfig := &aws.Config{
		Credentials: cred,
		Region:      aws.String("ap-northeast-1"),
	}

	config, err := readConf(params)
	if err != nil {
		log.Fatalln(err)
	}

	result, err := deploy.Deploy(awsConfig, config.Setting)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
	return 0
}

func (c *Deploy) Synopsis() string {
	synopsis := ""
	synopsis += "usage: thor deploy [options ...]\n"
	synopsis += "options:\n"
	synopsis += "    -f thor_setting.json\n"
	synopsis += "\n"
	synopsis += "    --profile=${aws profile name}\n"
	synopsis += "        --profile option is arbitrary parameter.\n"
	synopsis += "    --vault-password-file=${vault pass file}"
	synopsis += "\n"
	synopsis += "    --ask-vault-pass=${vault pass string}\n"
	synopsis += "===================================================\n"

	return synopsis
}
