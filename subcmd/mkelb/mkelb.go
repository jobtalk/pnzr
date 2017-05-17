package mkelb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/jobtalk/pnzr/api"
	"github.com/jobtalk/pnzr/lib/setting"
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

type mkelbConfigure struct {
	*setting.Setting
}

type mkelbParam struct {
	File    *string
	Profile *string
}

func (p mkelbParam) String() string {
	bin, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(bin)
}

type MkELB struct{}

func parseMkELBArgs(args []string) (*mkelbParam, error) {
	var result = &mkelbParam{}
	fileParam, err := getValFromArgs(args, "-f")
	if err != nil {
		return nil, err
	} else if len(fileParam) >= 2 {
		return nil, errors.New("'-f' parameter can not be specified more than once.")
	}
	if len(fileParam) == 1 {
		result.File = fileParam[0]
	} else if len(fileParam) == 0 {
		return nil, errors.New("'-f' parameter is a required item.")
	}
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

func (c *MkELB) Help() string {
	help := ""
	help += "usage: mkelb [options ...]\n"
	help += "options:\n"
	help += "    -f pnzr_setting.json\n"
	help += "\n"
	help += "    --profile=${aws profile name}\n"
	help += "        --profile option is arbitrary parameter.\n"

	return help
}

func (c *MkELB) Run(args []string) int {
	var config = &mkelbConfigure{}
	params, err := parseMkELBArgs(args)
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
	mkelbConfigureJSON, err := ioutil.ReadFile(*params.File)
	if err != nil {
		log.Fatalln(err)
	}
	err = json.Unmarshal(mkelbConfigureJSON, config)
	if err != nil {
		log.Fatalln(err)
	}
	result, err := api.MkELB(awsConfig, config.Setting.ELB)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
	return 0
}

func (c *MkELB) Synopsis() string {
	synopsis := ""
	synopsis += "usage: pnzr mkelb [options ...]\n"
	synopsis += "options:\n"
	synopsis += "    -f pnzr_setting.json\n"
	synopsis += "\n"
	synopsis += "    --profile=${aws profile name}\n"
	synopsis += "        --profile option is arbitrary parameter.\n"
	synopsis += "===================================================\n"

	return synopsis
}
