package v1

import (
	"github.com/jobtalk/pnzr/lib/setting"
	"path/filepath"
	"os"
	"errors"
	"regexp"
	"io/ioutil"
	"github.com/cbroglie/mustache"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var re = regexp.MustCompile(`.*\.json$`)

var (
	FileNotFoundError = errors.New("file info is nil")
)

type v1Setting struct {
	Version        float64
	Service        *ecs.CreateServiceInput
	TaskDefinition *ecs.RegisterTaskDefinitionInput
}

func (s *v1Setting) Convert() *setting.Setting {
	return &setting.Setting{
		s.Version,
		s.Service,
		s.TaskDefinition,
	}
}

func fileList(root string) ([]string, error) {
	if root == "" {
		return nil, nil
	}
	ret := []string{}
	err := filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if info == nil {
				return FileNotFoundError
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

type SettingLoader struct {
	sess     *session.Session
	kmsKeyID *string
}

func NewLoader(sess *session.Session, kmsKeyID *string) *SettingLoader {
	return &SettingLoader{
		sess:     sess,
		kmsKeyID: kmsKeyID,
	}
}

func (s *SettingLoader) Load(basePath, varsPath, outerVals string) (*setting.Setting, error) {
	var ret = v1Setting{}
	valueFileNameList, err := fileList(varsPath)
	if err != nil {
		return nil, err
	}

	templateFile, err := ioutil.ReadFile(basePath)
	if err != nil {
		return nil, err
	}

	for _, valueFileName := range valueFileNameList {
		var values = map[string]string{}
		valueBin, err := ioutil.ReadFile(valueFileName)
		if err != nil {
			return nil, err
		}

		if s.isEncrypt(valueBin) {

		}

		if err := json.Unmarshal(valueBin, &values); err != nil {
			return nil, err
		}

		result, err := mustache.Render(string(templateFile), values)
		if err != nil {
			return nil, err
		}

		templateFile = []byte(result)
	}

	if err := json.Unmarshal(templateFile, &ret); err != nil {
		return nil, err
	}
	return ret.Convert(), nil
}

func (s *SettingLoader) isEncrypt(bin []byte) bool {
	return false
}

func (s *SettingLoader) decrypt(bin []byte) ([]byte, error) {
	return nil, nil
}
