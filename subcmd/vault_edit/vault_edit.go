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

func (v *VaultEdit) parseArgs(args []string) {
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

func (v *VaultEdit) encrypt(keyID string, fileName string) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMS(v.sess)
	_, err = kms.SetKeyID(keyID).Encrypt(bin)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, []byte(kms.String()), 0644)
}

func (v *VaultEdit) decrypt(keyID string, fileName string) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMSFromBinary(bin, v.sess)
	if kms == nil {
		return errors.New(fmt.Sprintf("%v form is illegal", fileName))
	}
	plainText, err := kms.SetKeyID(keyID).Decrypt()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, plainText, 0644)
}

type VaultEdit struct {
	sess *session.Session
}

func (c *VaultEdit) Help() string {
	var msg string
	msg += "usage: pnzr vault-edit [options ...]\n"
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
	if e := os.Getenv("PNZR_EDITOR"); e != "" {
		return e
	}

	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}

	return "nano"
}

func (v *VaultEdit) Run(args []string) int {
	v.parseArgs(args)
	if *f == "" && *file == "" && len(flagSet.Args()) != 0 {
		targetName := flagSet.Args()[0]
		file = &targetName
	}

	if *file == "" {
		file = f
	}

	if err := v.decrypt(*kmsKeyID, *file); err != nil {
		log.Fatalln(err)
	}

	cmd := exec.Command(getEditor(), *file)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	if err := v.encrypt(*kmsKeyID, *file); err != nil {
		log.Fatalln(err)
	}

	return 0
}
