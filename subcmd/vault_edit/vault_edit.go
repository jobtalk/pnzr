package vedit

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/jobtalk/eriri/lib"
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

func encrypt(keyID string, fileName string, awsConfig *aws.Config) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMS()
	_, err = kms.SetKeyID(keyID).SetAWSConfig(awsConfig).Encrypt(bin)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, []byte(kms.String()), 0644)
}

func decrypt(keyID string, fileName string, awsConfig *aws.Config) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMSFromBinary(bin)
	if kms == nil {
		return errors.New(fmt.Sprintf("%v form is illegal", fileName))
	}
	plainText, err := kms.SetKeyID(keyID).SetAWSConfig(awsConfig).Decrypt()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, plainText, 0644)
}

type VaultEdit struct{}

func (c *VaultEdit) Help() string {
	var msg string
	msg += "usage: eriri vault-edit [options ...]\n"
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

func (c *VaultEdit) Synopsis() string {
	return c.Help()
}

func getEditor() string {
	if e := os.Getenv("eriri_EDITOR"); e != "" {
		return e
	}

	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}

	return "nano"
}

func (c *VaultEdit) Run(args []string) int {
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

	if err := decrypt(*kmsKeyID, *file, awsConfig); err != nil {
		log.Fatalln(err)
	}

	cmd := exec.Command(getEditor(), *file)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	if err := encrypt(*kmsKeyID, *file, awsConfig); err != nil {
		log.Fatalln(err)
	}

	return 0
}
