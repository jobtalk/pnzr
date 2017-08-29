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

type VaultViewCommand struct {
	sess           *session.Session
	kmsKeyID       *string
	file           *string
	f              *string
	profile        *string
	region         *string
	awsAccessKeyID *string
	awsSecretKeyID *string
	flagSet        *flag.FlagSet
}

func (v *VaultViewCommand) parseArgs(args []string) {
	v.kmsKeyID = v.flagSet.String("key_id", getenv.String("KMS_KEY_ID"), "Amazon KMS key ID")
	v.profile = v.flagSet.String("profile", getenv.String("AWS_PROFILE_NAME", "default"), "aws credentials profile name")
	v.region = v.flagSet.String("region", getenv.String("AWS_REGION", "ap-northeast-1"), "aws region")
	v.awsAccessKeyID = v.flagSet.String("aws-access-key-id", getenv.String("AWS_ACCESS_KEY_ID"), "aws access key id")
	v.awsSecretKeyID = v.flagSet.String("aws-secret-key-id", getenv.String("AWS_SECRET_KEY_ID"), "aws secret key id")
	v.file = v.flagSet.String("file", "", "target file")
	v.f = v.flagSet.String("f", "", "target file")

	if err := v.flagSet.Parse(args); err != nil {
		log.Fatalln(err)
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
}

func (v *VaultViewCommand) decrypt(keyID string, fileName string, awsConfig *aws.Config) ([]byte, error) {
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

func (c *VaultViewCommand) Help() string {
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

func (c *VaultViewCommand) Synopsis() string {
	return c.Help()
}

func (v *VaultViewCommand) Run(args []string) int {
	v.flagSet = new(flag.FlagSet)
	v.parseArgs(args)

	if *v.f == "" && *v.file == "" && len(v.flagSet.Args()) != 0 {
		targetName := v.flagSet.Args()[0]
		v.file = &targetName
	}

	var cred *credentials.Credentials
	if *v.awsAccessKeyID != "" && *v.awsSecretKeyID != "" {
		cred = credentials.NewStaticCredentials(*v.awsAccessKeyID, *v.awsSecretKeyID, "")
	} else {
		cred = credentials.NewSharedCredentials("", *v.profile)
	}

	awsConfig := &aws.Config{
		Credentials: cred,
		Region:      v.region,
	}

	if *v.file == "" {
		v.file = v.f
	}

	plain, err := v.decrypt(*v.kmsKeyID, *v.file, awsConfig)
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
