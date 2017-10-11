package v1

import (
	"encoding/json"
	"github.com/ieee0824/cryptex"
	"github.com/jobtalk/pnzr/vars"
	"io/ioutil"
	"reflect"
)

var (
	ORIGIN_TEST_DIR = vars.TEST_DATA_DIR_ROOT + "/subcmd/vault/encrypt/v1"
	COPY_TEST_DIR   = ORIGIN_TEST_DIR + "/copy"
	KEY_DIR         = vars.TEST_DATA_DIR_ROOT + "/key"
)

type decrypter struct {
	crypter *cryptex.Cryptex
}

func (d *decrypter) decrypt(kmsKeyID, fileName string) error {
	container := new(cryptex.Container)

	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bin, container); err != nil {
		return err
	}

	i, err := d.crypter.Decrypt(container)
	if err != nil {
		return err
	}
	plain, err := json.MarshalIndent(i, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, plain, 0644)
}
