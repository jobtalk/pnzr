package v1

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/cryptex"
	"github.com/ieee0824/cryptex/kms"
	"io/ioutil"
)

func Encrypt(sess *session.Session, keyID, fileName string) error {
	var plain interface{}
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bin, &plain); err != nil {
		return err
	}
	kmsClient := kms.New(sess)
	kmsClient.SetKey(keyID)
	c := cryptex.New(kmsClient)

	cipher, err := c.Encrypt(plain)
	if err != nil {
		return err
	}

	chipherBin, err := json.MarshalIndent(cipher, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, chipherBin, 0644)
}
