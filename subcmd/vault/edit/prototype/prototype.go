package prototype

import (
	"io/ioutil"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/jobtalk/pnzr/lib"
	"errors"
	"os"
	"os/exec"
	"github.com/jobtalk/pnzr/subcmd/vault/edit/util"
)



func encrypt(sess *session.Session, keyID, fileName string) error {
	bin, err := ioutil.ReadFile(fileName)
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

func decrypt(sess *session.Session, fileName string) error {
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

func Edit(sess *session.Session, keyID, fileName string) error {
	if err := decrypt(sess, fileName); err != nil {
		panic(err)
	}
	defer func() {
		if err := encrypt(sess, keyID, fileName); err != nil {
			panic(err)
		}
	}()
	cmd := exec.Command(util.GetEditor(), fileName)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		panic(err)
	}
	return nil
}