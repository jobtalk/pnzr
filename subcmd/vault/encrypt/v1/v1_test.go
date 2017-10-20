package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"

	"github.com/ieee0824/cryptex"
	"github.com/ieee0824/cryptex/rsa"
	"github.com/jobtalk/pnzr/vars"
)

type encrypterTester struct {
	ORIGIN_TEST_DIR string
	COPY_TEST_DIR   string
	KEY_DIR         string
	privateKey      []byte
	publicKey       []byte
	cryptex         *cryptex.Cryptex
	encrypter       *Encrypter
	decrypter       *decrypter
}

func newEncrypterTester() (*encrypterTester, error) {
	var ret = new(encrypterTester)

	ret.ORIGIN_TEST_DIR = vars.TEST_DATA_DIR_ROOT + "/subcmd/vault/encrypt/v1"
	ret.COPY_TEST_DIR = ret.ORIGIN_TEST_DIR + "/copy"
	ret.KEY_DIR = vars.TEST_DATA_DIR_ROOT + "/key"

	if err := ret.generateRsaKey(); err != nil {
		return nil, err
	}

	ret.generateCryptex()
	ret.generateEncrypter()
	ret.generateDecrypter()

	if err := ret.generateTestFiles(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *encrypterTester) generateRsaKey() error {
	var err error
	t.privateKey, err = ioutil.ReadFile(t.KEY_DIR + "/private-key.pem")
	if err != nil {
		return err
	}
	t.publicKey, err = ioutil.ReadFile(t.KEY_DIR + "/public-key.pem")
	if err != nil {
		return err
	}
	return nil
}

func (t *encrypterTester) generateCryptex() {
	t.cryptex = cryptex.New(rsa.New(t.privateKey, t.publicKey))
}

func (t *encrypterTester) generateEncrypter() {
	t.encrypter = &Encrypter{t.cryptex}
}

func (t *encrypterTester) generateDecrypter() {
	t.decrypter = &decrypter{t.cryptex}
}

func (t *encrypterTester) generateTestFiles() error {
	if err := exec.Command("mkdir", "-p", t.COPY_TEST_DIR).Run(); err != nil {
		return err
	}
	infos, err := ioutil.ReadDir(t.ORIGIN_TEST_DIR)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.IsDir() {
			continue
		}
		src, err := os.Open(t.ORIGIN_TEST_DIR + "/" + info.Name())
		if err != nil {
			panic(err)
		}
		dst, err := os.Create(t.COPY_TEST_DIR + "/" + info.Name())
		if err != nil {
			panic(err)
		}
		if _, err := io.Copy(dst, src); err != nil {
			panic(err)
		}
	}

	return nil
}

func (t *encrypterTester) fin() error {
	if err := os.RemoveAll(t.COPY_TEST_DIR); err != nil {
		return err
	}
	return nil
}

func (tester *encrypterTester) run(t *testing.T, tests map[string]testInput, f func(t *testing.T, tester *encrypterTester, tests map[string]testInput) error) error {
	return f(t, tester, tests)
}

func encryptTestRun(t *testing.T, tester *encrypterTester, tests map[string]testInput) error {
	for key, test := range tests {
		err := tester.encrypter.Encrypt("", fmt.Sprintf("%s/%v", tester.COPY_TEST_DIR, key))
		if err != nil {
			return err
		}

		if !test.expectsErr {
			org, err := ioutil.ReadFile(tester.ORIGIN_TEST_DIR + "/" + key)
			if err != nil {
				return err
			}
			encrypted, err := ioutil.ReadFile(tester.COPY_TEST_DIR + "/" + key)
			if err != nil {
				return err
			}
			if compareJSON(org, encrypted) {
				t.Log(string(org))
				t.Log(string(encrypted))
				t.Fatalf("equal")
			}

			if err := tester.decrypter.decrypt(tester.COPY_TEST_DIR + "/" + key); err != nil {
				return err
			}

			decrypted, err := ioutil.ReadFile(tester.COPY_TEST_DIR + "/" + key)
			if !test.expectsErr && err != nil {
				t.Fatalf("should not be error but: %v", err)
			}
			if test.expectsErr && err == nil {
				t.Fatalf("should be error")
			}
			if !compareJSON(org, decrypted) {
				t.Log(string(org))
				t.Log(string(decrypted))
				t.Fatalf("not match")
			}
		}
	}
	return nil
}

func alreadyEncryptErrorTestRun(t *testing.T, tester *encrypterTester, tests map[string]testInput) error {
	for key, _ := range tests {
		err := tester.encrypter.Encrypt("", fmt.Sprintf("%s/%v", tester.COPY_TEST_DIR, key))
		if err != nil {
			return err
		}
		err = tester.encrypter.Encrypt("", fmt.Sprintf("%s/%v", tester.COPY_TEST_DIR, key))
		if err == nil {
			t.Fatalf("should be error")
		}
	}
	return nil
}

func compareJSON(a, b []byte) bool {
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

type testInput struct {
	expectsErr bool
}

func TestEncrypter_Encrypt(t *testing.T) {
	tests := map[string]testInput{
		"0.json": {
			false,
		},
	}
	tester, err := newEncrypterTester()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tester.fin(); err != nil {
			panic(err)
		}
	}()

	if err := tester.run(t, tests, encryptTestRun); err != nil {
		panic(err)
	}
}

func TestEncrypter_AlreadyEncrypt(t *testing.T) {
	tests := map[string]testInput{
		"0.json": {
			true,
		},
	}
	tester, err := newEncrypterTester()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := tester.fin(); err != nil {
			panic(err)
		}
	}()

	if err := tester.run(t, tests, alreadyEncryptErrorTestRun); err != nil {
		panic(err)
	}
}