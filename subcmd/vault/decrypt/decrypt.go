package decrypt

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/subcmd/vault/decrypt/prototype"
	"github.com/jobtalk/pnzr/subcmd/vault/decrypt/v1"
	"github.com/jobtalk/pnzr/vars"
)

type DecryptCommand struct {
	sess           *session.Session
	file           *string
	profile        *string
	region         *string
	awsAccessKeyID *string
	awsSecretKeyID *string
	configVersion  *string
}

func (d *DecryptCommand) parseArgs(args []string) (helpString string) {
	var (
		flagSet = new(flag.FlagSet)
		f       *string
	)

	buffer := new(bytes.Buffer)
	flagSet.SetOutput(buffer)
	d.profile = flagSet.String("profile", getenv.String("AWS_PROFILE_NAME", "default"), "aws credentials profile name")
	d.region = flagSet.String("region", getenv.String("AWS_REGION", "ap-northeast-1"), "aws region")
	d.awsAccessKeyID = flagSet.String("aws-access-key-id", getenv.String("AWS_ACCESS_KEY_ID"), "aws access key id")
	d.awsSecretKeyID = flagSet.String("aws-secret-key-id", getenv.String("AWS_SECRET_KEY_ID"), "aws secret key id")
	d.file = flagSet.String("file", "", "target file")
	d.configVersion = flagSet.String("v", vars.CONFIG_VERSION, "config version")
	f = flagSet.String("f", "", "target file")

	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return buffer.String()
		}
		panic(err)
	}

	if *f == "" && *d.file == "" && len(flagSet.Args()) != 0 {
		targetName := flagSet.Args()[0]
		d.file = &targetName
	}

	if *d.file == "" {
		d.file = f
	}

	var awsConfig = aws.Config{}

	if *d.awsAccessKeyID != "" && *d.awsSecretKeyID != "" && *d.profile == "" {
		awsConfig.Credentials = credentials.NewStaticCredentials(*d.awsAccessKeyID, *d.awsSecretKeyID, "")
		awsConfig.Region = d.region
	}

	d.sess = session.Must(session.NewSessionWithOptions(session.Options{
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
		Profile:                 *d.profile,
		Config:                  awsConfig,
	}))

	return
}

func (d *DecryptCommand) decrypt(fileName string) error {
	switch *d.configVersion {
	case "1.0":
		return v1.Decrypt(d.sess, fileName)
	case "prototype":
		return prototype.Decrypt(d.sess, fileName)
	default:
		return fmt.Errorf("unsupport configure version: %v", *d.configVersion)
	}
}

func (d *DecryptCommand) Help() string {
	return d.parseArgs([]string{"-h"})
}

func (d *DecryptCommand) Synopsis() string {
	return "Decryption mode of encrypted file."
}

func (d *DecryptCommand) Run(args []string) int {
	if len(args) == 0 {
		fmt.Println(d.Synopsis())
		return 0
	}
	d.parseArgs(args)

	if err := d.decrypt(*d.file); err != nil {
		panic(err)
	}

	return 0
}
