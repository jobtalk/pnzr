package vault

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/lib"
)

type VaultCommand struct {
	sess           *session.Session
	kmsKeyID       *string
	encryptFlag    *bool
	decryptFlag    *bool
	file           *string
	profile        *string
	region         *string
	awsAccessKeyID *string
	awsSecretKeyID *string
}

func (v *VaultCommand) parseArgs(args []string) {
	var (
		flagSet = new(flag.FlagSet)
		f       *string
	)

	v.kmsKeyID = flagSet.String("key_id", getenv.String("KMS_KEY_ID"), "Amazon KMS key ID")
	v.encryptFlag = flagSet.Bool("encrypt", getenv.Bool("ENCRYPT", false), "encrypt mode")
	v.decryptFlag = flagSet.Bool("decrypt", getenv.Bool("DECRYPT", false), "decrypt mode")
	v.profile = flagSet.String("profile", getenv.String("AWS_PROFILE_NAME", "default"), "aws credentials profile name")
	v.region = flagSet.String("region", getenv.String("AWS_REGION", "ap-northeast-1"), "aws region")
	v.awsAccessKeyID = flagSet.String("aws-access-key-id", getenv.String("AWS_ACCESS_KEY_ID"), "aws access key id")
	v.awsSecretKeyID = flagSet.String("aws-secret-key-id", getenv.String("AWS_SECRET_KEY_ID"), "aws secret key id")
	v.file = flagSet.String("file", "", "target file")
	f = flagSet.String("f", "", "target file")

	if err := flagSet.Parse(args); err != nil {
		log.Fatalln(err)
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
}

func (v *VaultCommand) encrypt(keyID string, fileName string) error {
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

func (v *VaultCommand) decrypt(keyID string, fileName string) error {
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

func (c *VaultCommand) Help() string {
	var msg string
	msg += "usage: pnzr vault [options ...]\n"
	msg += "options:\n"
	msg += "    -key_id\n"
	msg += "        set kms key id\n"
	msg += "    -encrypt\n"
	msg += "        use encrypt mode\n"
	msg += "    -decrypt\n"
	msg += "        use decrypt mode\n"
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

func (v *VaultCommand) Run(args []string) int {
	v.parseArgs(args)

	if *v.encryptFlag == *v.decryptFlag {
		log.Fatalln("Choose whether to execute encrypt or decrypt.")
	}
	if *v.decryptFlag {
		err := v.decrypt(*v.kmsKeyID, *v.file)
		if err != nil {
			log.Fatalln(err)
		}
	} else if *v.encryptFlag {
		err := v.encrypt(*v.kmsKeyID, *v.file)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return 0
}

func (c *VaultCommand) Synopsis() string {
	return c.Help()
}
