package deploy

import (
	"bytes"
	"flag"

	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/api"
	"github.com/jobtalk/pnzr/lib/config"
	"github.com/jobtalk/pnzr/lib/config/v0"
	"github.com/jobtalk/pnzr/lib/config/v1"
	"strings"
)

func parseDockerImage(image string) (url, tag string) {
	r := strings.Split(image, ":")
	if len(r) == 2 {
		return r[0], r[1]
	}
	return r[0], ""
}

type DeployCommand struct {
	sess           *session.Session
	file           *string
	profile        *string
	kmsKeyID       *string
	region         *string
	externalPath   *string
	outerVals      *string
	awsAccessKeyID *string
	awsSecretKeyID *string
	tagOverride    *string
}

func (d *DeployCommand) parseArgs(args []string) (helpString string) {
	flagSet := new(flag.FlagSet)
	var f *string

	buffer := new(bytes.Buffer)
	flagSet.SetOutput(buffer)

	d.kmsKeyID = flagSet.String("key_id", getenv.String("KMS_KEY_ID"), "Amazon KMS key ID")
	d.file = flagSet.String("file", "", "target file")
	f = flagSet.String("f", "", "target file")
	d.profile = flagSet.String("profile", getenv.String("AWS_PROFILE_NAME", "default"), "aws credentials profile name")
	d.region = flagSet.String("region", getenv.String("AWS_REGION", "ap-northeast-1"), "aws region")
	d.externalPath = flagSet.String("vars_path", getenv.String("PNZR_VARS_PATH"), "external conf path")
	d.outerVals = flagSet.String("V", "", "outer values")
	d.tagOverride = flagSet.String("t", getenv.String("DOCKER_DEFAULT_DEPLOY_TAG", "latest"), "tag override param")
	d.awsAccessKeyID = flagSet.String("aws-access-key-id", getenv.String("AWS_ACCESS_KEY_ID"), "aws access key id")
	d.awsSecretKeyID = flagSet.String("aws-secret-key-id", getenv.String("AWS_SECRET_KEY_ID"), "aws secret key id")

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

func (d *DeployCommand) Run(args []string) int {
	d.parseArgs(args)
	var conf *config.IntermediateConfig
	if v1_config.CheckSupportVersion(d.file) {
		loader := &v1_config.ConfigLoader{}
		c, err := loader.Load(d.file)
		if err != nil {
			panic(err)
		}
		conf = c

	} else {
		loader := v0_config.NewLoader(d.sess, d.kmsKeyID)

		c, err := loader.Load(*d.file, *d.externalPath, *d.outerVals)
		if err != nil {
			panic(err)
		}

		conf = c
	}

	for i, containerDefinition := range conf.TaskDefinition.ContainerDefinitions {
		imageName, tag := parseDockerImage(*containerDefinition.Image)
		if tag == "$tag" {
			image := imageName + ":" + *d.tagOverride
			conf.TaskDefinition.ContainerDefinitions[i].Image = &image
		} else if tag == "" {
			image := imageName + ":" + "latest"
			conf.TaskDefinition.ContainerDefinitions[i].Image = &image
		}
	}

	result, err := api.Deploy(d.sess, conf)
	if err != nil {
		panic(err)
	}

	resultJSON, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(resultJSON))

	return 0
}

func (c *DeployCommand) Synopsis() string {
	return "Deploy docker on ecs."
}

func (c *DeployCommand) Help() string {
	return c.parseArgs([]string{"-h"})
}
