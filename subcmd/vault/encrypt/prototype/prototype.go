package prototype

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jobtalk/pnzr/lib"
	"io/ioutil"
)

func Encrypt(sess *session.Session, keyID, fileName string) error {
	bin, err :=ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMS(sess)
	_, err = kms.SetKeyID(keyID).Encrypt(bin)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(fileName, []byte(kms.String()), 0644)
}
