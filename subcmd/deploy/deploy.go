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
	"strconv"
	"strings"
	"time"

	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/api"
	"github.com/jobtalk/pnzr/lib"
	"github.com/jobtalk/pnzr/lib/setting"
)

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
	dryRun         *bool
	progress       *bool
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

func parseDockerImage(image string) (url, tag string) {
	r := strings.Split(image, ":")
	if len(r) == 2 {
		return r[0], r[1]
	}
	return r[0], ""
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
	plainText, err := kms.SetKeyID(*d.kmsKeyID).Decrypt()
	if err != nil {
		return nil, err
	}
	return plainText, nil
}

func (d *DeployCommand) readConf(base []byte, externalPathList []string) (*deployConfigure, error) {
	var root = *d.externalPath
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
	d.dryRun = flagSet.Bool("dry-run", false, "dry run mode")
	d.progress = flagSet.Bool("progress", false, "show progress")

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
	var config = &deployConfigure{}

	externalList, err := fileList(*d.externalPath)
	if err != nil {
		log.Fatalln(err)
	}
	baseConfBinary, err := ioutil.ReadFile(*d.file)
	if err != nil {
		log.Fatal(err)
	}

	if *d.outerVals != "" {
		baseStr, err := lib.Embedde(string(baseConfBinary), *d.outerVals)
		if err == nil {
			baseConfBinary = []byte(baseStr)
		}
	}

	if externalList != nil {
		c, err := d.readConf(baseConfBinary, externalList)
		if err != nil {
			log.Fatalln(err)
		}
		config = c
	} else {
		bin, err := ioutil.ReadFile(*d.file)
		if err != nil {
			log.Fatalln(err)
		}
		if err := json.Unmarshal(bin, config); err != nil {
			log.Fatalln(err)
		}
	}

	for i, containerDefinition := range config.ECS.TaskDefinition.ContainerDefinitions {
		imageName, tag := parseDockerImage(*containerDefinition.Image)
		if tag == "$tag" {
			image := imageName + ":" + *d.tagOverride
			config.ECS.TaskDefinition.ContainerDefinitions[i].Image = &image
		} else if tag == "" {
			image := imageName + ":" + "latest"
			config.ECS.TaskDefinition.ContainerDefinitions[i].Image = &image
		}
	}

	if *d.dryRun {
		dryRunFormat := &DryRun{
			*d.region,
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
	t := result.([]interface{})[0]
	id := *t.(*ecs.RegisterTaskDefinitionOutput).TaskDefinition.Revision

	if *d.progress {
		input := &ecs.DescribeServicesInput{
			Services: []*string{
				config.ECS.Service.ServiceName,
			},
			Cluster: config.ECS.Service.Cluster,
		}
		s := ecs.New(d.sess)
		fmt.Printf("(1/3) 【%s】の【%s】へのデプロイ を開始\n", *config.ECS.Service.Cluster, *config.ECS.Service.ServiceName)
		ok := false
		flagNext := false
		for !ok {
			t, err := s.DescribeServices(input)
			if err != nil {
				panic(err)
			}
			deployments := t.Services[0].Deployments
			s := strings.Split(*deployments[0].TaskDefinition, ":")
			nextID, err := strconv.Atoi(s[len(s)-1])
			if err != nil {
				panic(err)
			}
			if int(id) == nextID && len(deployments) > 1 && *deployments[0].DesiredCount == *deployments[0].RunningCount && !flagNext {
				fmt.Printf("(2/3) 【%s】の【%s】へのデプロイは新しいコンテナを起動\n", *config.ECS.Service.Cluster, *config.ECS.Service.ServiceName)
				flagNext = true
			}
			if int(id) == nextID && flagNext && len(deployments) == 1 {
				fmt.Printf("(3/3) 【%s】の【%s】へのデプロイは古いコンテナの停止\n", *config.ECS.Service.Cluster, *config.ECS.Service.ServiceName)
				ok = true
			}
			time.Sleep(3 * time.Second)
		}
		fmt.Printf("【%s】の【%s】へのデプロイが終了\n", *config.ECS.Service.Cluster, *config.ECS.Service.ServiceName)
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
