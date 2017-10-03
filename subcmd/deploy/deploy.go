package deploy

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/api"
	"github.com/jobtalk/pnzr/lib/setting"
	"github.com/jobtalk/pnzr/lib/setting/prototype"
	"github.com/jobtalk/pnzr/lib/setting/v1"
)

type DeployCommand struct {
	dryRun             *bool
	sess               *session.Session
	config             *aws.Config
	credentialFileName string
	paramsFromArgs     *params
	paramsFromEnvs     *params
	mergedParams       *params
}

type params struct {
	kmsKeyID     *string
	file         *string
	profile      *string
	region       *string
	varsPath     *string
	overrideTag  *string
	awsAccessKey *string
	awsSecretKey *string
}

type DryRun struct {
	Region string
	Config setting.Setting
}

func (d DryRun) String() string {
	structJSON, err := json.MarshalIndent(d, "", "   ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", string(structJSON))
}

var (
	DockerImageParseErr       = errors.New("parse error")
	IllegalAccessKeyOptionErr = errors.New("There was an illegal input in '-aws-access-key-id' or '-aws-secret-key-id '")
)

func parseDockerImage(image string) (url, tag string, err error) {
	r := strings.Split(image, ":")
	if 3 <= len(r) {
		return "", "", DockerImageParseErr
	}
	if len(r) == 2 {
		return r[0], r[1], nil
	}
	return r[0], "", nil
}

func stringIsEmpty(s *string) bool {
	if s == nil {
		return true
	}

	return len(*s) == 0
}

func (d *DeployCommand) parseArgs(args []string) (helpString string) {
	p := params{}
	flagSet := new(flag.FlagSet)
	var f *string

	buffer := new(bytes.Buffer)
	flagSet.SetOutput(buffer)

	p.kmsKeyID = flagSet.String("key_id", "", "Amazon KMS key ID")
	p.file = flagSet.String("file", "", "target file")
	f = flagSet.String("f", "", "target file")
	p.profile = flagSet.String("profile", "", "aws credentials profile name")
	p.region = flagSet.String("region", "", "aws region")
	p.varsPath = flagSet.String("vars_path", "", "external conf path")
	p.overrideTag = flagSet.String("t", "", "tag override param")
	p.awsAccessKey = flagSet.String("aws-access-key-id", "", "aws access key id")
	p.awsSecretKey = flagSet.String("aws-secret-key-id", "", "aws secret key id")
	d.dryRun = flagSet.Bool("dry-run", false, "dry run mode")

	if err := flagSet.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return buffer.String()
		}
		panic(err)
	}

	if *f == "" && *p.file == "" && len(flagSet.Args()) != 0 {
		targetName := flagSet.Args()[0]
		p.file = &targetName
	}

	if *p.file == "" {
		p.file = f
	}

	d.paramsFromArgs = &p

	return
}

func (d *DeployCommand) parseEnv() {
	p := params{}

	p.kmsKeyID = aws.String(getenv.String("KMS_KEY_ID"))
	p.profile = aws.String(getenv.String("AWS_PROFILE_NAME", "default"))
	p.overrideTag = aws.String(getenv.String("DOCKER_DEFAULT_DEPLOY_TAG", "latest"))
	p.region = aws.String(getenv.String("AWS_REGION"))
	p.awsAccessKey = aws.String(getenv.String("AWS_ACCESS_KEY_ID"))
	p.awsSecretKey = aws.String(getenv.String("AWS_SECRET_ACCESS_KEY"))

	d.paramsFromEnvs = &p
}

func (d *DeployCommand) mergeParams() {
	result := params{}

	if d.paramsFromArgs == nil {
		d.mergedParams = d.paramsFromEnvs
		return
	}
	if d.paramsFromEnvs == nil {
		d.mergedParams = d.paramsFromArgs
		return
	}

	if stringIsEmpty(d.paramsFromArgs.kmsKeyID) && !stringIsEmpty(d.paramsFromEnvs.kmsKeyID) {
		result.kmsKeyID = d.paramsFromEnvs.kmsKeyID
	} else {
		result.kmsKeyID = d.paramsFromArgs.kmsKeyID
	}

	if stringIsEmpty(d.paramsFromArgs.file) && !stringIsEmpty(d.paramsFromEnvs.file) {
		result.file = d.paramsFromEnvs.file
	} else {
		result.file = d.paramsFromArgs.file
	}

	if stringIsEmpty(d.paramsFromArgs.varsPath) && !stringIsEmpty(d.paramsFromEnvs.varsPath) {
		result.varsPath = d.paramsFromEnvs.varsPath
	} else {
		result.varsPath = d.paramsFromArgs.varsPath
	}

	if stringIsEmpty(d.paramsFromArgs.profile) && !stringIsEmpty(d.paramsFromEnvs.profile) {
		result.profile = d.paramsFromEnvs.profile
	} else {
		result.profile = d.paramsFromArgs.profile
	}

	if stringIsEmpty(d.paramsFromArgs.overrideTag) && !stringIsEmpty(d.paramsFromEnvs.overrideTag) {
		result.overrideTag = d.paramsFromEnvs.overrideTag
	} else {
		result.overrideTag = d.paramsFromArgs.overrideTag
	}

	if stringIsEmpty(d.paramsFromArgs.region) && !stringIsEmpty(d.paramsFromEnvs.region) {
		result.region = d.paramsFromEnvs.region
	} else {
		result.region = d.paramsFromArgs.region
	}

	if stringIsEmpty(d.paramsFromArgs.awsAccessKey) && !stringIsEmpty(d.paramsFromEnvs.awsAccessKey) {
		result.awsAccessKey = d.paramsFromEnvs.awsAccessKey
	} else {
		result.awsAccessKey = d.paramsFromArgs.awsAccessKey
	}

	if stringIsEmpty(d.paramsFromArgs.awsSecretKey) && !stringIsEmpty(d.paramsFromEnvs.awsSecretKey) {
		result.awsSecretKey = d.paramsFromEnvs.awsSecretKey
	} else {
		result.awsSecretKey = d.paramsFromArgs.awsSecretKey
	}

	d.mergedParams = &result
}

func (d *DeployCommand) loadCredentials() error {
	if stringIsEmpty(d.mergedParams.awsAccessKey) != stringIsEmpty(d.mergedParams.awsSecretKey) {
		return IllegalAccessKeyOptionErr
	}

	if !stringIsEmpty(d.mergedParams.awsAccessKey) && !stringIsEmpty(d.mergedParams.awsSecretKey) {
		d.config = new(aws.Config)
		d.config.Credentials = credentials.NewStaticCredentials(*d.mergedParams.awsAccessKey, *d.mergedParams.awsSecretKey, "")
		return nil
	}

	return nil
}

func (d *DeployCommand) generateSession() {
	if d.config != nil {
		d.sess = session.Must(session.NewSessionWithOptions(session.Options{
			AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
			SharedConfigState:       session.SharedConfigEnable,
			Config:                  *d.config,
		}))
	} else if !stringIsEmpty(d.mergedParams.profile) {
		d.sess = session.Must(session.NewSessionWithOptions(session.Options{
			AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
			SharedConfigState:       session.SharedConfigEnable,
			Profile:                 *d.mergedParams.profile,
		}))
	} else {
		d.sess = session.Must(session.NewSessionWithOptions(session.Options{
			AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
			SharedConfigState:       session.SharedConfigEnable,
			Profile:                 "default",
		}))
	}

	if !stringIsEmpty(d.mergedParams.region) {
		d.sess.Config.Region = d.mergedParams.region
	}
}

func (d *DeployCommand) Run(args []string) int {
	var deploySetting = &setting.Setting{}
	var loader setting.Loader
	d.parseArgs(args)
	d.parseEnv()
	d.mergeParams()
	if err := d.loadCredentials(); err != nil {
		panic(err)
	}
	d.generateSession()

	if setting.IsV1Setting(*d.mergedParams.file) {
		loader = v1.NewLoader(d.sess, d.mergedParams.kmsKeyID)
	} else {
		loader = prototype.NewLoader(d.sess, d.mergedParams.kmsKeyID)
	}

	s, err := loader.Load(*d.mergedParams.file, *d.mergedParams.varsPath, "")
	if err != nil {
		panic(err)
	}
	deploySetting = s

	for i, containerDefinition := range deploySetting.TaskDefinition.ContainerDefinitions {
		imageName, tag, err := parseDockerImage(*containerDefinition.Image)
		if err != nil {
			panic(err)
		}

		if tag == "$tag" {
			image := imageName + ":" + *d.mergedParams.overrideTag
			deploySetting.TaskDefinition.ContainerDefinitions[i].Image = &image
		} else if tag == "" {
			image := imageName + ":" + "latest"
			deploySetting.TaskDefinition.ContainerDefinitions[i].Image = &image
		}
	}

	if *d.dryRun {
		dryRunFormat := &DryRun{
			*d.mergedParams.region,
			*deploySetting,
		}
		f, err := os.Open("/dev/stderr")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(f, "******** DRY RUN ********\n%s\n", *dryRunFormat)
		f.Close()
		return 0
	}

	result, err := api.Deploy(d.sess, deploySetting)
	if err != nil {
		log.Fatalln(err)
	}
	resultJSON, err := json.MarshalIndent(result, "", "    ")
	if err != nil {
		log.Fatalln(err)
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
