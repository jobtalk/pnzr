package v1

import (
	"encoding/json"
	"fmt"
	"github.com/ieee0824/cryptex"
	"github.com/ieee0824/cryptex/rsa"
	"github.com/jobtalk/pnzr/vars"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

var (
	ORIGIN_TEST_DIR = vars.TEST_DATA_DIR_ROOT + "/subcmd/vault/encrypt/v1"
	COPY_TEST_DIR   = ORIGIN_TEST_DIR + "/copy"
	KEY_DIR         = vars.TEST_DATA_DIR_ROOT + "/key"
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

type decrypter struct {
	crypter *cryptex.Cryptex
}

func (d *decrypter) decrypt(fileName string) error {
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

func TestEncrypter_Encrypt(t *testing.T) {
	privateKey, err := ioutil.ReadFile(KEY_DIR + "/private-key.pem")
	if err != nil {
		panic(err)
	}
	publicKey, err := ioutil.ReadFile(KEY_DIR + "/public-key.pem")
	if err != nil {
		panic(err)
	}

	crypter := cryptex.New(rsa.New(privateKey, publicKey))
	encrypter := &Encrypter{
		crypter,
	}
	testDecrypter := &decrypter{crypter}

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

	tests := map[string]struct {
		err bool
	}{
		"0.json": {
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
	}

	for key, test := range tests {
		err := encrypter.Encrypt("", fmt.Sprintf("%s/%v", COPY_TEST_DIR, key))
		if err != nil {
			panic(err)
		}

		if !test.err {
			org, err := ioutil.ReadFile(ORIGIN_TEST_DIR + "/" + key)
			if err != nil {
				panic(err)
			}
			encrypted, err := ioutil.ReadFile(COPY_TEST_DIR + "/" + key)
			if err != nil {
				panic(err)
			}
			if compaireJSON(org, encrypted) {
				t.Log(string(org))
				t.Log(string(encrypted))
				t.Fatalf("equal")
			}

			if err := testDecrypter.decrypt(COPY_TEST_DIR + "/" + key); err != nil {
				panic(err)
			}

			decrypted, err := ioutil.ReadFile(COPY_TEST_DIR + "/" + key)
			if !test.err && err != nil {
				t.Fatalf("should not be error but: %v", err)
			}
			if test.err && err == nil {
				t.Fatalf("should be error")
			}
			if !compaireJSON(org, decrypted) {
				t.Log(string(org))
				t.Log(string(decrypted))
				t.Fatalf("not match")
			}
		}
	}
}
