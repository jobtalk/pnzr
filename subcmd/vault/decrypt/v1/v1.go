package v1

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/cryptex"
	"github.com/ieee0824/cryptex/kms"
	"io/ioutil"
)

func Decrypt(sess *session.Session, fileName string) error {
	var chipher = &cryptex.Container{}
	chipherBin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(chipherBin, &chipher); err != nil {
		return err
	}

	c := cryptex.New(kms.New(sess))

	plain, err := c.Decrypt(chipher)
	if err != nil {
		return err
	}
	plainBin, err := json.MarshalIndent(plain, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, plainBin, 0644)
}
