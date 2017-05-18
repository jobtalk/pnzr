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

func init() {
	kmsKeyID = flagSet.String("key_id", "", "Amazon KMS key ID")
	profile = flagSet.String("profile", "default", "aws credentials profile name")
	region = flagSet.String("region", "ap-northeast-1", "aws region")

	awsAccessKeyID = flagSet.String("aws-access-key-id", "", "aws access key id")
	awsSecretKeyID = flagSet.String("aws-secret-key-id", "", "aws secret key id")

	file = flagSet.String("file", "", "target file")
	f = flagSet.String("f", "", "target file")
}

func decrypt(keyID string, fileName string, awsConfig *aws.Config) ([]byte, error) {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	kms := lib.NewKMSFromBinary(bin)
	if kms == nil {
		return nil, errors.New(fmt.Sprintf("%v form is illegal", fileName))
	}
	plainText, err := kms.SetKeyID(keyID).SetAWSConfig(awsConfig).Decrypt()
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

type VaultView struct{}

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

func (c *VaultView) Run(args []string) int {
	if err := flagSet.Parse(args); err != nil {
		log.Fatalln(err)
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

	plain, err := decrypt(*kmsKeyID, *file, awsConfig)
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
