package encrypt

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/subcmd/vault/encrypt/prototype"
	"github.com/jobtalk/pnzr/subcmd/vault/encrypt/v1"
	"github.com/jobtalk/pnzr/subcmd/vault/encrypt/v1/iface"
	"github.com/jobtalk/pnzr/vars"
)

type EncryptCommand struct {
	sess           *session.Session
	kmsKeyID       *string
	file           *string
	profile        *string
	region         *string
	awsAccessKeyID *string
	awsSecretKeyID *string
	configVersion  *string
	v1Encrypter    v1_api.API
}

func (e *EncryptCommand) encrypt(keyID string, fileName string) error {
	switch *e.configVersion {
	case "1.0":
		return e.v1Encrypter.Encrypt(keyID, fileName)
	case "prototype":
		return prototype.Encrypt(e.sess, keyID, fileName)
	default:
		return fmt.Errorf("unsupport configure version: %v", *e.configVersion)
	}
}

func (e *EncryptCommand) parseArgs(args []string) (helpString string) {
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

func (e *EncryptCommand) Help() string {
	return e.parseArgs([]string{"-h"})
}

func (e *EncryptCommand) Synopsis() string {
	return "Encryption mode of plaintext file."

}

func (e *EncryptCommand) Run(args []string) int {
	if len(args) == 0 {
		fmt.Println(e.Synopsis())
		return 0
	}
	e.parseArgs(args)
	e.v1Encrypter = v1.New(e.sess, *e.kmsKeyID)

	if err := e.encrypt(*e.kmsKeyID, *e.file); err != nil {
		panic(err)
	}

	return 0
}
