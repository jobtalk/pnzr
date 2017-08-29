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

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/getenv"
	"github.com/jobtalk/pnzr/api"
	"github.com/jobtalk/pnzr/lib"
	"github.com/jobtalk/pnzr/lib/setting"
)

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

func (d *DeployCommand) parseArgs(args []string) {
	flagSet := new(flag.FlagSet)
	var f *string

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
		log.Fatalln(err)
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
	synopsis := ""
	synopsis += "usage: pnzr deploy [options ...]\n"
	synopsis += "options:\n"
	synopsis += "    -f pnzr_setting.json\n"
	synopsis += "\n"
	synopsis += "    -profile=${aws profile name}\n"
	synopsis += "        -profile option is arbitrary parameter.\n"
	synopsis += "    -region\n"
	synopsis += "        aws region\n"
	synopsis += "    -vars_path\n"
	synopsis += "        setting external values path file\n"
	synopsis += "    -V\n"
	synopsis += "        setting outer values\n"
	synopsis += "    -aws-access-key-id\n"
	synopsis += "        setting aws access key id\n"
	synopsis += "    -aws-secret-key-id\n"
	synopsis += "        setting aws secret key id\n"
	synopsis += "    -t tag name\n"
	synopsis += "        setting docker tag\n"
	synopsis += "        defining as follows will replace $tag.\n"
	synopsis += "        \"Image\":\"image-name:$tag\"\n"
	synopsis += "    -key_id \n"
	synopsis += "        set kms key id\n"
	synopsis += "===================================================\n"

	return synopsis
}

func (c *DeployCommand) Help() string {
	return c.Synopsis()
}
