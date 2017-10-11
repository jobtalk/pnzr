package v1

import (
	"testing"
	"github.com/jobtalk/pnzr/vars"
	"os/exec"
	"os"
	"github.com/ieee0824/cryptex"
	"io/ioutil"
	"github.com/ieee0824/cryptex/rsa"
	"fmt"
	"io"
	"encoding/json"
	"reflect"
)

var (
	ORIGIN_TEST_DIR = vars.TEST_DATA_DIR_ROOT + "/subcmd/vault/decrypt/v1"
	COPY_TEST_DIR = ORIGIN_TEST_DIR + "/copy"
	KEY_DIR = vars.TEST_DATA_DIR_ROOT + "/key"
)

func compaireJSON(a, b []byte) bool {
	var ai, bi interface{}
	if err := json.Unmarshal(a, &ai); err != nil {
		panic(err)
	}

	if err := json.Unmarshal(b, &bi); err != nil {
		panic(err)
	}

	return reflect.DeepEqual(ai, bi)
}

type encrypter struct {
	crypter *cryptex.Cryptex
}

func (e *encrypter) encrypt(fileName string) error {
	var i interface{}
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bin, &i); err != nil {
		return err
	}

	c, err := e.crypter.Encrypt(i)
	if err != nil {
		return err
	}

	chipher, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, chipher, 0644)
}

func TestDecrypter_Decrypt(t *testing.T) {
	privateKey, err := ioutil.ReadFile(KEY_DIR + "/private-key.pem")
	if err != nil {
		panic(err)
	}
	publicKey, err := ioutil.ReadFile(KEY_DIR + "/public-key.pem")
	if err != nil {
		panic(err)
	}

	crypter := cryptex.New(rsa.New(privateKey, publicKey))
	decrypter := &Decrypter{
		crypter,
	}
	testEncrypter := &encrypter{crypter}

	if err := exec.Command("mkdir", "-p", COPY_TEST_DIR).Run(); err != nil {
		panic(err)
	}
	defer func() {
		if err := os.RemoveAll(COPY_TEST_DIR); err != nil {
			panic(err)
		}
	}()

	infos, err := ioutil.ReadDir(ORIGIN_TEST_DIR)
	if err != nil {
		panic(err)
	}


	tests := map[string]struct{
		isEncrypt bool
		err bool
	}{
		"0.json": {
			false,
			true,
		},
		"1.json": {
			true,
			false,
		},
	}

	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		src, err := os.Open(ORIGIN_TEST_DIR + "/" + info.Name())
		if err != nil {
			panic(err)
		}
		dst, err := os.Create(COPY_TEST_DIR + "/" + info.Name())
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(dst, src); err != nil {
			panic(err)
		}

		test, ok := tests[info.Name()]
		if ok && test.isEncrypt {
			err := testEncrypter.encrypt(COPY_TEST_DIR + "/" + info.Name())
			if err != nil {
				panic(err)
			}
		}
	}

	for i, test := range tests {
		err := decrypter.Decrypt(fmt.Sprintf("%s/%v", COPY_TEST_DIR, i))
		if !test.err && err != nil {
			t.Fatalf("should not be error but: %v", err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error")
		}
		if !test.err {
			org, err := ioutil.ReadFile(ORIGIN_TEST_DIR + "/" + i)
			if err != nil {
				panic(err)
			}
			decrypted, err := ioutil.ReadFile(COPY_TEST_DIR + "/" + i)
			if err != nil {
				panic(err)
			}
			if !compaireJSON(org, decrypted) {
				t.Log(string(org))
				t.Log(string(decrypted))
				t.Fatalf("not match")
			}
		}
	}
}