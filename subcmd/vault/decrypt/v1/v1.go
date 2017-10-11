package v1

import (
	"encoding/json"
	"github.com/ieee0824/cryptex"
	"github.com/ieee0824/cryptex/kms"
	"io/ioutil"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Decrypter struct {
	c *cryptex.Cryptex
}

func New(s *session.Session) *Decrypter {
	return &Decrypter{
		cryptex.New(kms.New(s)),
	}
}

func (d *Decrypter)Decrypt(fileName string) error {
	var chipher = &cryptex.Container{}
	chipherBin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(chipherBin, &chipher); err != nil {
		return err
	}

	plain, err := d.c.Decrypt(chipher)
	if err != nil {
		return err
	}
	plainBin, err := json.MarshalIndent(plain, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, plainBin, 0644)
}
