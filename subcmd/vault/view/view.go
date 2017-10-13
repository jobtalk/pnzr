package view

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/cryptex"
	"github.com/ieee0824/cryptex/kms"
	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/lib"
	"github.com/jobtalk/pnzr/vars"
	"io/ioutil"
	"os"
	"os/exec"
)

type ViewCommand struct {
	sess           *session.Session
	file           *string
	profile        *string
	region         *string
	awsAccessKeyID *string
	awsSecretKeyID *string
	configVersion  *string
}

func (v *ViewCommand) parseArgs(args []string) (helpString string) {
	var (
		flagSet = new(flag.FlagSet)
		f       *string
	)

	buffer := new(bytes.Buffer)
	flagSet.SetOutput(buffer)

	v.profile = flagSet.String("profile", getenv.String("AWS_PROFILE_NAME", "default"), "aws credentials profile name")
	v.region = flagSet.String("region", getenv.String("AWS_REGION", "ap-northeast-1"), "aws region")
	v.awsAccessKeyID = flagSet.String("aws-access-key-id", getenv.String("AWS_ACCESS_KEY_ID"), "aws access key id")
	v.awsSecretKeyID = flagSet.String("aws-secret-key-id", getenv.String("AWS_SECRET_KEY_ID"), "aws secret key id")
	v.file = flagSet.String("file", "", "target file")
	v.configVersion = flagSet.String("v", vars.CONFIG_VERSION, "config version")
	f = flagSet.String("f", "", "target file")

	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return buffer.String()
		}
		panic(err)
	}

	if *f == "" && *v.file == "" && len(flagSet.Args()) != 0 {
		targetName := flagSet.Args()[0]
		v.file = &targetName
	}

	if *v.file == "" {
		v.file = f
	}

	var awsConfig = aws.Config{}

	if *v.awsAccessKeyID != "" && *v.awsSecretKeyID != "" && *v.profile == "" {
		awsConfig.Credentials = credentials.NewStaticCredentials(*v.awsAccessKeyID, *v.awsSecretKeyID, "")
		awsConfig.Region = v.region
	}

	v.sess = session.Must(session.NewSessionWithOptions(session.Options{
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
		Profile:                 *v.profile,
		Config:                  awsConfig,
	}))

	return
}

func (v *ViewCommand) decryptTemporary(fileName string) ([]byte, error) {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	kms := lib.NewKMSFromBinary(bin, v.sess)
	if kms == nil {
		return nil, errors.New(fmt.Sprintf("%v form is illegal", fileName))
	}
	plainText, err := kms.Decrypt()
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

func (v *ViewCommand) Help() string {
	return v.parseArgs([]string{"-h"})
}

func (v *ViewCommand) Synopsis() string {
	return "Viewer mode of encrypted file."
}

func (v *ViewCommand) Run(args []string) int {
	if len(args) == 0 {
		fmt.Println(v.Synopsis())
		return 0
	}
	v.parseArgs(args)

	switch *v.configVersion {
	case "1.0":
		var container = &cryptex.Container{}
		bin, err := ioutil.ReadFile(*v.file)
		if err != nil {
			panic(err)
		}
		if err := json.Unmarshal(bin, container); err != nil {
			panic(err)
		}
		if err := cryptex.New(kms.New(v.sess)).View(container); err != nil {
			panic(err)
		}
	case "prototype":
		plain, err := v.decryptTemporary(*v.file)
		if err != nil {
			panic(err)
		}

		cmd := exec.Command("less")
		cmd.Stdin = bytes.NewReader(plain)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			panic(err)
		}
	default:
		panic(fmt.Errorf("unsupport version: %v", *v.configVersion))
	}

	return 0
}
