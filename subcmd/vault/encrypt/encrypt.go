package encrypt

import (
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/lib"
	"io/ioutil"
)

type EncryptCommand struct {
	sess           *session.Session
	kmsKeyID       *string
	file           *string
	profile        *string
	region         *string
	awsAccessKeyID *string
	awsSecretKeyID *string
}

func (e *EncryptCommand) encrypt(keyID string, fileName string) error {
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

func (e *EncryptCommand) parseArgs(args []string) {
	var (
		flagSet = new(flag.FlagSet)
		f       *string
	)

	e.kmsKeyID = flagSet.String("key_id", getenv.String("KMS_KEY_ID"), "Amazon KMS key ID")
	e.profile = flagSet.String("profile", getenv.String("AWS_PROFILE_NAME", "default"), "aws credentials profile name")
	e.region = flagSet.String("region", getenv.String("AWS_REGION", "ap-northeast-1"), "aws region")
	e.awsAccessKeyID = flagSet.String("aws-access-key-id", getenv.String("AWS_ACCESS_KEY_ID"), "aws access key id")
	e.awsSecretKeyID = flagSet.String("aws-secret-key-id", getenv.String("AWS_SECRET_KEY_ID"), "aws secret key id")
	e.file = flagSet.String("file", "", "target file")
	f = flagSet.String("f", "", "target file")

	if err := flagSet.Parse(args); err != nil {
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
}

func (e *EncryptCommand) Help() string {
	return ""
}

func (e *EncryptCommand) Synopsis() string {
	return e.Help()
}

func (e *EncryptCommand) Run(args []string) int {
	e.parseArgs(args)

	if err := e.encrypt(*e.kmsKeyID, *e.file); err != nil {
		panic(err)
	}

	return 0
}
