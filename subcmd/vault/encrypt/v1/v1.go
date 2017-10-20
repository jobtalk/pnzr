package v1

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/cryptex"
	"github.com/ieee0824/cryptex/kms"
	"io/ioutil"
	"fmt"
)

type Encrypter struct {
	c *cryptex.Cryptex
}

func New(s *session.Session, keyID string) *Encrypter {
	return &Encrypter{
		cryptex.New(kms.New(s).SetKey(keyID)),
	}
}

func (e *Encrypter) Encrypt(keyID, fileName string) error {
	var plain map[string]interface{}
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bin, &plain); err != nil {
		return err
	}

	if _, ok := plain["encryption_type"]; ok {
		return fmt.Errorf("%s is already encrypted", fileName)
	}

	cipher, err := e.c.Encrypt(plain)
	if err != nil {
		return err
	}

	chipherBin, err := json.MarshalIndent(cipher, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, chipherBin, 0644)
}
