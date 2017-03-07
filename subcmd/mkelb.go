package subcmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/ieee0824/thor/mkelb"
	"github.com/ieee0824/thor/setting"
)

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
	help += "    -f thor_setting.json\n"
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
	result, err := mkelb.MkELB(awsConfig, config.Setting.ELB)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(result)
	return 0
}

func (c *MkELB) Synopsis() string {
	synopsis := ""
	synopsis += "usage: thor mkelb [options ...]\n"
	synopsis += "options:\n"
	synopsis += "    -f thor_setting.json\n"
	synopsis += "\n"
	synopsis += "    --profile=${aws profile name}\n"
	synopsis += "        --profile option is arbitrary parameter.\n"
	synopsis += "===================================================\n"

	return synopsis
}
