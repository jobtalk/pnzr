package deploy

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/api"
	"github.com/jobtalk/pnzr/lib"
	"github.com/jobtalk/pnzr/lib/setting"
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
	ECS    setting.ECS
}

func (d DryRun) String() string {
	structJSON, err := json.MarshalIndent(d, "", "   ")
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s", string(structJSON))
}

var re = regexp.MustCompile(`.*\.json$`)

var (
	DockerImageParseErr       = errors.New("parse error")
	IllegalAccessKeyOptionErr = errors.New("There was an illegal input in '-aws-access-key-id' or '-aws-secret-key-id '")
)

func stringIsEmpty(s *string) bool {
	if s == nil {
		return true
	}

	return len(*s) == 0
}

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

func fileList(root string) ([]string, error) {
	if root == "" {
		return nil, nil
	}
	ret := []string{}
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return errors.New("file info is nil")
			}
			if info.IsDir() {
				return nil
			}

			rel, err := filepath.Rel(root, path)
			if re.MatchString(rel) {
				ret = append(ret, rel)
			}

			return nil
		})

	if err != nil {
		return nil, err
	}

	return ret, nil
}

type deployConfigure struct {
	*setting.Setting
}

func isEncrypted(data []byte) bool {
	var buffer = map[string]interface{}{}
	if err := json.Unmarshal(data, &buffer); err != nil {
		return false
	}
	elem, ok := buffer["cipher"]
	if !ok {
		return false
	}
	str, ok := elem.(string)
	if !ok {
		return false
	}

	return len(str) != 0
}

func (d *DeployCommand) decrypt(bin []byte) ([]byte, error) {
	kms := lib.NewKMSFromBinary(bin, d.sess)
	if kms == nil {
		return nil, errors.New(fmt.Sprintf("%v format is illegal", string(bin)))
	}
	plainText, err := kms.SetKeyID(*d.mergedParams.kmsKeyID).Decrypt()
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

func (d *DeployCommand) readConf(base []byte, externalPathList []string) (*deployConfigure, error) {
	var root = *d.mergedParams.varsPath
	var ret = &deployConfigure{}
	baseStr := string(base)

	root = strings.TrimSuffix(root, "/")
	for _, externalPath := range externalPathList {
		external, err := ioutil.ReadFile(root + "/" + externalPath)
		if err != nil {
			return nil, err
		}
		if isEncrypted(external) {
			plain, err := d.decrypt(external)
			if err != nil {
				return nil, err
			}
			external = plain
		}
		baseStr, err = lib.Embedde(baseStr, string(external))
		if err != nil {
			return nil, err
		}
	}
	if err := json.Unmarshal([]byte(baseStr), ret); err != nil {
		return nil, err
	}
	return ret, nil
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
	d.parseArgs(args)
	d.parseEnv()
	d.mergeParams()
	if err := d.loadCredentials(); err != nil {
		panic(err)
	}
	d.generateSession()
	var config = &deployConfigure{}

	externalList, err := fileList(*d.mergedParams.varsPath)
	if err != nil {
		log.Fatalln(err)
	}
	baseConfBinary, err := ioutil.ReadFile(*d.mergedParams.file)
	if err != nil {
		log.Fatal(err)
	}

	if externalList != nil {
		c, err := d.readConf(baseConfBinary, externalList)
		if err != nil {
			log.Fatalln(err)
		}
		config = c
	} else {
		bin, err := ioutil.ReadFile(*d.mergedParams.file)
		if err != nil {
			log.Fatalln(err)
		}
		if err := json.Unmarshal(bin, config); err != nil {
			log.Fatalln(err)
		}
	}

	for i, containerDefinition := range config.ECS.TaskDefinition.ContainerDefinitions {
		imageName, tag, err := parseDockerImage(*containerDefinition.Image)
		if err != nil {
			panic(err)
		}
		if tag == "$tag" {
			image := imageName + ":" + *d.mergedParams.overrideTag
			config.ECS.TaskDefinition.ContainerDefinitions[i].Image = &image
		} else if tag == "" {
			image := imageName + ":" + "latest"
			config.ECS.TaskDefinition.ContainerDefinitions[i].Image = &image
		}
	}

	if *d.dryRun {
		dryRunFormat := &DryRun{
			*d.mergedParams.region,
			*config.ECS,
		}
		f, err := os.Open("/dev/stderr")
		if err != nil {
			panic(err)
		}
		fmt.Fprintf(f, "******** DRY RUN ********\n%s\n", *dryRunFormat)
		f.Close()
		return 0
	}

	result, err := api.Deploy(d.sess, config.Setting)
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
