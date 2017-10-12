package edit

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/lib"
	"github.com/jobtalk/pnzr/vars"
	"io/ioutil"
	"github.com/jobtalk/pnzr/subcmd/vault/edit/prototype"
	"github.com/jobtalk/pnzr/subcmd/vault/edit/v1/iface"
	"github.com/jobtalk/pnzr/subcmd/vault/edit/v1"
)



type EditCommand struct {
	sess           *session.Session
	kmsKeyID       *string
	file           *string
	profile        *string
	region         *string
	awsAccessKeyID *string
	awsSecretKeyID *string
	configVersion  *string
	v1Editor v1_api.API
}

func (e *EditCommand) parseArgs(args []string) (helpString string) {
	var (
		flagSet = new(flag.FlagSet)
		f       *string
	)

	buffer := new(bytes.Buffer)
	flagSet.SetOutput(buffer)

	e.kmsKeyID = flagSet.String("key_id", getenv.String("KMS_KEY_ID"), "Amazon KMS key ID")
	e.profile = flagSet.String("profile", getenv.String("AWS_PROFILE_NAME", "default"), "aws credentials profile name")
	e.region = flagSet.String("region", getenv.String("AWS_REGION", "ap-northeast-1"), "aws region")
	e.awsAccessKeyID = flagSet.String("aws-access-key-id", getenv.String("AWS_ACCESS_KEY_ID"), "aws access key id")
	e.awsSecretKeyID = flagSet.String("aws-secret-key-id", getenv.String("AWS_SECRET_KEY_ID"), "aws secret key id")
	e.file = flagSet.String("file", "", "target file")
	e.configVersion = flagSet.String("v", vars.CONFIG_VERSION, "config version")
	f = flagSet.String("f", "", "target file")

	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return buffer.String()
		}
		panic(err)
	}

	if *f == "" && *e.file == "" && len(flagSet.Args()) != 0 {
		targetName := flagSet.Args()[0]
		e.file = &targetName
	}

	if *e.file == "" {
		e.file = f
	}

	var awsConfig = aws.Config{}

	if *e.awsAccessKeyID != "" && *e.awsSecretKeyID != "" && *e.profile == "" {
		awsConfig.Credentials = credentials.NewStaticCredentials(*e.awsAccessKeyID, *e.awsSecretKeyID, "")
		awsConfig.Region = e.region
	}

	e.sess = session.Must(session.NewSessionWithOptions(session.Options{
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
		Profile:                 *e.profile,
		Config:                  awsConfig,
	}))

	return
}

func (e *EditCommand) decrypt(fileName string) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMSFromBinary(bin, e.sess)
	if kms == nil {
		return errors.New(fmt.Sprintf("%v form is illegal", fileName))
	}
	plainText, err := kms.Decrypt()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, plainText, 0644)
}

func (e *EditCommand) encrypt(keyID string, fileName string) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMS(e.sess)
	_, err = kms.SetKeyID(keyID).Encrypt(bin)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, []byte(kms.String()), 0644)
}

func (e *EditCommand) Help() string {
	return e.parseArgs([]string{"-h"})
}

func (e *EditCommand) Synopsis() string {
	return "Edit mode of encrypted file."
}

func (e *EditCommand) Run(args []string) int {
	if len(args) == 0 {
		fmt.Println(e.Synopsis())
		return 0
	}
	e.parseArgs(args)
	e.v1Editor = v1.New(e.sess, *e.kmsKeyID)

	switch *e.configVersion {
	case "1.0":
		if err := e.v1Editor.Edit(*e.file); err != nil {
			panic(err)
		}
	case "prototype":
		if err := prototype.Edit(e.sess, *e.kmsKeyID, *e.file); err != nil {
			panic(err)
		}
	default:
		panic(fmt.Errorf("unsupport version: %v", *e.configVersion))
	}

	return 0
}
