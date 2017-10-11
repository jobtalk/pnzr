package prototype

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"io/ioutil"
	"fmt"
	"github.com/jobtalk/pnzr/lib"
	"errors"
)

func Decrypt(sess *session.Session, fileName string) error {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	kms := lib.NewKMSFromBinary(bin, sess)
	if kms == nil {
		return errors.New(fmt.Sprintf("%v form is illegal", fileName))
	}
	plainText, err := kms.Decrypt()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, plainText, 0644)
}