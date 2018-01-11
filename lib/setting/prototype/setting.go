package prototype

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	intermediate "github.com/jobtalk/pnzr/lib/setting"
	"github.com/jobtalk/pnzr/lib/setting/prototype/embedde"
	"github.com/jobtalk/pnzr/lib/setting/prototype/kms"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`.*\.json$`)

var (
	BadReqKMS = errors.New("bad request: kms error")
)

type ECS struct {
	Service        *ecs.CreateServiceInput
	TaskDefinition *ecs.RegisterTaskDefinitionInput
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

type setting struct {
	ECS *ECS
}

func (s *setting) version() float64 {
	return 0.0
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

func (s *SettingLoader) Load(basePath, varsPath, outerVals string) (*intermediate.Setting, error) {
	varsFileList, err := fileList(varsPath)
	if err != nil {
		return nil, err
	}

	baseConfBinary, err := ioutil.ReadFile(basePath)
	if err != nil {
		return nil, err
	}

	if outerVals != "" {
		baseStr, err := embedde.Embedde(string(baseConfBinary), outerVals)
		if err == nil {
			baseConfBinary = []byte(baseStr)
		}
	}

	if len(varsFileList) != 0 {
		result, err := s.loadConf(baseConfBinary, varsPath, varsFileList)
		if err != nil {
			return nil, err
		}
		var ret = intermediate.Setting{}
		ret.Version = result.version()
		ret.Service = result.ECS.Service
		ret.TaskDefinition = result.ECS.TaskDefinition

		return &ret, nil
	}

	var result = &setting{}
	if err := json.Unmarshal(baseConfBinary, result); err != nil {
		return nil, err
	}

	var ret = intermediate.Setting{}
	ret.Version = result.version()
	ret.Service = result.ECS.Service
	ret.TaskDefinition = result.ECS.TaskDefinition
	return &ret, nil
}

func (s *SettingLoader) loadConf(base []byte, varsRoot string, varsFileNameList []string) (*setting, error) {
	var (
		ret     = &setting{}
		baseStr = string(base)
	)
	varsRoot = strings.TrimSuffix(varsRoot, "/")

	for _, varsFileName := range varsFileNameList {
		varsBinary, err := ioutil.ReadFile(fmt.Sprintf("%s/%s", varsRoot, varsFileName))
		if err != nil {
			return nil, err
		}

		if s.isEncrypted(varsBinary) {
			plain, err := s.decrypt(varsBinary)
			if err != nil {
				return nil, err
			}

			varsBinary = plain
		}
		baseStr, err = embedde.Embedde(baseStr, string(varsBinary))
		if err != nil {
			return nil, err
		}
	}

	if err := json.Unmarshal([]byte(baseStr), ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (*SettingLoader) isEncrypted(data []byte) bool {
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

func (s *SettingLoader) decrypt(bin []byte) ([]byte, error) {
	k := kms.NewKMSFromBinary(bin, s.sess)
	if k == nil {
		return nil, BadReqKMS
	}

	plainText, err := k.SetKeyID(*s.kmsKeyID).Decrypt()
	if err != nil {
		return nil, err
	}
	return plainText, nil
}
