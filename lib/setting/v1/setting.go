package v1

import (
	"encoding/json"
	"errors"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/ieee0824/cryptex"
	"github.com/ieee0824/cryptex/kms"
	"github.com/ieee0824/jec"
	"github.com/jobtalk/pnzr/lib/setting"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	FileNotFoundError = errors.New("file info is nil")
)

var re = regexp.MustCompile(`.*\.json$`)

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
	reg      []*regexp.Regexp
}

func NewLoader(sess *session.Session, kmsKeyID *string) *SettingLoader {
	return &SettingLoader{
		sess:     sess,
		kmsKeyID: kmsKeyID,
	}
}

func (s *SettingLoader) Load(basePath, varsPath, outerVals string) (*setting.Setting, error) {
	var setting = v1Setting{}
	varFileList, err := fileList(varsPath)
	if err != nil {
		return nil, err
	}
	baseConf, err := ioutil.ReadFile(basePath)
	if err != nil {
		return nil, err
	}

	for _, path := range varFileList {
		bin, err := ioutil.ReadFile(strings.TrimSuffix(varsPath, "/") + "/" + path)
		if err != nil {
			return nil, err
		}

		if b, err := s.decrypt(bin); err != nil {
			return nil, err
		} else {
			bin = b
		}

		baseConf, err = jec.Embed(baseConf, bin)
		if err != nil {
			return nil, err
		}
	}

	if err := json.Unmarshal(baseConf, &setting); err != nil {
		return nil, err
	}

	return setting.Convert(), nil
}

func (s *SettingLoader) isEncrypt(bin []byte) bool {
	var buffer = cryptex.Container{}
	if err := json.Unmarshal(bin, &buffer); err != nil {
		return false
	}
	return buffer.EncryptionType == "kms"
}

func (s *SettingLoader) decrypt(bin []byte) ([]byte, error) {
	if !s.isEncrypt(bin) {
		return bin, nil
	}

	kmsClient := kms.New(s.sess)
	kmsClient.SetKey(*s.kmsKeyID)

	var buffer = cryptex.Container{}
	if err := json.Unmarshal(bin, &buffer); err != nil {
		return nil, err
	}
	plain, err := cryptex.New(kmsClient).Decrypt(&buffer)
	if err != nil {
		return nil, err
	}
	return json.Marshal(plain)
}
