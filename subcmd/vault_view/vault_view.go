package vview

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/lib"
)

var flagSet = &flag.FlagSet{}

var (
	kmsKeyID       *string
	file           *string
	f              *string
	profile        *string
	region         *string
	awsAccessKeyID *string
	awsSecretKeyID *string
)

func (v *VaultView) parseArgs(args []string) {
	kmsKeyID = flagSet.String("key_id", getenv.String("KMS_KEY_ID"), "Amazon KMS key ID")
	profile = flagSet.String("profile", getenv.String("AWS_PROFILE_NAME", "default"), "aws credentials profile name")
	region = flagSet.String("region", getenv.String("AWS_REGION", "ap-northeast-1"), "aws region")

	awsAccessKeyID = flagSet.String("aws-access-key-id", getenv.String("AWS_ACCESS_KEY_ID"), "aws access key id")
	awsSecretKeyID = flagSet.String("aws-secret-key-id", getenv.String("AWS_SECRET_KEY_ID"), "aws secret key id")

	file = flagSet.String("file", "", "target file")
	f = flagSet.String("f", "", "target file")

	if err := flagSet.Parse(args); err != nil {
		log.Fatalln(err)
	}

	var awsConfig = aws.Config{}

	if *awsAccessKeyID != "" && *awsSecretKeyID != "" && *profile == "" {
		awsConfig.Credentials = credentials.NewStaticCredentials(*awsAccessKeyID, *awsSecretKeyID, "")
		awsConfig.Region = region
	}

	v.sess = session.Must(session.NewSessionWithOptions(session.Options{
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
		Profile:                 *profile,
		Config:                  awsConfig,
	}))
}

func (v *VaultView) decrypt(keyID string, fileName string, awsConfig *aws.Config) ([]byte, error) {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	kms := lib.NewKMSFromBinary(bin, v.sess)
	if kms == nil {
		return nil, errors.New(fmt.Sprintf("%v form is illegal", fileName))
	}
	plainText, err := kms.SetKeyID(keyID).Decrypt()
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

type VaultView struct {
	sess *session.Session
}

func (c *VaultView) Help() string {
	var msg string
	msg += "usage: pnzr vault-view [options ...]\n"
	msg += "options:\n"
	msg += "    -key_id\n"
	msg += "        set kms key id\n"
	msg += "    -file\n"
	msg += "        setting target file\n"
	msg += "    -f"
	msg += "        setting target file\n"
	msg += "    -profile\n"
	msg += "        aws credential name\n"
	msg += "    -region\n"
	msg += "        aws region name\n"
	msg += "    -aws-access-key-id\n"
	msg += "        setting aws access key id\n"
	msg += "    -aws-secret-key-id\n"
	msg += "        setting aws secret key id\n"
	msg += "===================================================\n"
	return msg
}

func (c *VaultView) Synopsis() string {
	return c.Help()
}

func (v *VaultView) Run(args []string) int {
	v.parseArgs(args)

	if *f == "" && *file == "" && len(flagSet.Args()) != 0 {
		targetName := flagSet.Args()[0]
		file = &targetName
	}

	var cred *credentials.Credentials
	if *awsAccessKeyID != "" && *awsSecretKeyID != "" {
		cred = credentials.NewStaticCredentials(*awsAccessKeyID, *awsSecretKeyID, "")
	} else {
		cred = credentials.NewSharedCredentials("", *profile)
	}

	awsConfig := &aws.Config{
		Credentials: cred,
		Region:      region,
	}

	if *file == "" {
		file = f
	}

	plain, err := v.decrypt(*kmsKeyID, *file, awsConfig)
	if err != nil {
		log.Fatalln(err)
	}

	cmd := exec.Command("less")
	cmd.Stdin = bytes.NewReader(plain)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	return 0
}
