package edit

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/getenv"
	config "github.com/jobtalk/pnzr/lib/config/v1"
)

func getEditor() string {
	if e := os.Getenv("PNZR_EDITOR"); e != "" {
		return e
	}

	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}

	return "nano"
}

type EditCommand struct {
	sess           *session.Session
	kmsKeyID       *string
	file           *string
	profile        *string
	region         *string
	awsAccessKeyID *string
	awsSecretKeyID *string
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
	kms := config.NewKMSFromBinary(bin, e.sess)
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
	kms := config.NewKMS(e.sess)
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

	if err := e.decrypt(*e.file); err != nil {
		panic(err)
	}
	defer func() {
		if err := e.encrypt(*e.kmsKeyID, *e.file); err != nil {
			panic(err)
		}
	}()
	cmd := exec.Command(getEditor(), *e.file)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
	return 0
}
