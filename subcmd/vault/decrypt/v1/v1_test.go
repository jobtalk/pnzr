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

type decrypterTester struct {
	ORIGIN_TEST_DIR string
	COPY_TEST_DIR   string
	KEY_DIR         string
	privateKey      []byte
	publicKey       []byte
	cryptex         *cryptex.Cryptex
	encrypter       *encrypter
	decrypter       *Decrypter
	testFileInfos   []os.FileInfo
}

func newDecrypterTester() (*decrypterTester, error) {
	ret := &decrypterTester{}
	ret.ORIGIN_TEST_DIR = vars.TEST_DATA_DIR_ROOT + "/subcmd/vault/decrypt/v1"
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

func (t *decrypterTester) generateRsaKey() error {
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

func (t *decrypterTester) generateCryptex() {
	t.cryptex = cryptex.New(rsa.New(t.privateKey, t.publicKey))
}

func (t *decrypterTester) generateEncrypter() {
	t.encrypter = &encrypter{t.cryptex}
}

func (t *decrypterTester) generateDecrypter() {
	t.decrypter = &Decrypter{t.cryptex}
}

func (t *decrypterTester) generateTestFiles() error {
	if err := exec.Command("mkdir", "-p", t.COPY_TEST_DIR).Run(); err != nil {
		return err
	}
	infos, err := ioutil.ReadDir(t.ORIGIN_TEST_DIR)
	if err != nil {
		return err
	}
	t.testFileInfos = infos

	return nil
}

func (t *decrypterTester) generateEncryptedFiles(tests map[string]testInput) error {
	for _, info := range t.testFileInfos {
		if info.IsDir() {
			continue
		}
		src, err := os.Open(t.ORIGIN_TEST_DIR + "/" + info.Name())
		if err != nil {
			return err
		}
		dst, err := os.Create(t.COPY_TEST_DIR + "/" + info.Name())
		if err != nil {
			return err
		}
		if _, err := io.Copy(dst, src); err != nil {
			return err
		}

		test, ok := tests[info.Name()]
		if ok && test.encryptsTestFileAtFirst {
			err := t.encrypter.encrypt(t.COPY_TEST_DIR + "/" + info.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (t *decrypterTester) fin() error {
	if err := os.RemoveAll(t.COPY_TEST_DIR); err != nil {
		return err
	}
	return nil
}

func (tester *decrypterTester) run(t *testing.T, tests map[string]testInput) error {
	for i, test := range tests {
		err := tester.decrypter.Decrypt(fmt.Sprintf("%s/%v", tester.COPY_TEST_DIR, i))
		if !test.expectsErr && err != nil {
			t.Fatalf("should not be error but: %v", err)
		}
		if test.expectsErr && err == nil {
			t.Fatalf("should be error")
		}
		if !test.expectsErr {
			org, err := ioutil.ReadFile(tester.ORIGIN_TEST_DIR + "/" + i)
			if err != nil {
				return err
			}
			decrypted, err := ioutil.ReadFile(tester.COPY_TEST_DIR + "/" + i)
			if err != nil {
				return err
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

// testInput struct defines test data.
//
// encryptsTestFileAtFirst field is a flag
// that controls whether to encrypt at the start of the test.
//
// expectsErr field is a flag
// that controls whether or not to expect errors.
type testInput struct {
	encryptsTestFileAtFirst bool
	expectsErr              bool
}

func TestDecrypter_Decrypt(t *testing.T) {
	tester, err := newDecrypterTester()
	if err != nil {
		panic(err)
	}

	tests := map[string]testInput{
		"0.json": {
			false,
			true,
		},
		"1.json": {
			true,
			false,
		},
	}

	if err := tester.generateEncryptedFiles(tests); err != nil {
		panic(err)
	}

	defer func() {
		if err := tester.fin(); err != nil {
			panic(err)
		}
	}()

	if err := tester.run(t, tests); err != nil {
		panic(err)
	}
}
