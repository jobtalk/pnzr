package v1

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/cryptex"
	"github.com/ieee0824/cryptex/rsa"
	"github.com/jobtalk/pnzr/vars"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"
)

var (
	ORIGIN_TEST_DIR = vars.TEST_DATA_DIR_ROOT + "/subcmd/vault/edit/v1"
	COPY_TEST_DIR   = ORIGIN_TEST_DIR + "/copy"
	KEY_DIR         = vars.TEST_DATA_DIR_ROOT + "/key"
)

type mocChiper struct {
	fault bool
}

func newChiper(b bool) *mocChiper {
	return &mocChiper{b}
}

func (m *mocChiper) Encrypt(d []byte) ([]byte, error) {
	if m.fault {
		return nil, fmt.Errorf("some error")
	}
	d = append(d[1:], d[0])
	return d, nil
}

func (m *mocChiper) Decrypt(d []byte) ([]byte, error) {
	if m.fault {
		return nil, fmt.Errorf("some error")
	}
	d = append(d[len(d)-1:], d[:len(d)-1]...)
	return d, nil
}

func (m *mocChiper) EncryptionType() string {
	return "moc"
}

type encrypter struct {
	crypter *cryptex.Cryptex
}

func (d *encrypter) decrypt(fileName string) error {
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

func TestNew(t *testing.T) {
	if nil == New(session.New(), "") {
		t.Fatalf("not allocated")
	}
}

func TestEditor_Edit2(t *testing.T) {
	testEditor := &Editor{
		cryptex.New(newChiper(true)),
		fmt.Sprintf("%s/editor", COPY_TEST_DIR),
	}

	tests := []struct {
		fileName string
	}{
		{"hoge"},
		{ORIGIN_TEST_DIR + "/0.json"},
	}

	for _, test := range tests {
		err := testEditor.Edit(test.fileName)
		if err == nil {
			t.Fatalf("should be error for %v but not:", test.fileName)
		}
	}
}

func TestEditor_Edit1(t *testing.T) {
	privateKey, err := ioutil.ReadFile(KEY_DIR + "/private-key.pem")
	if err != nil {
		panic(err)
	}
	publicKey, err := ioutil.ReadFile(KEY_DIR + "/public-key.pem")
	if err != nil {
		panic(err)
	}

	if err := exec.Command("mkdir", "-p", COPY_TEST_DIR).Run(); err != nil {
		panic(err)
	}
	defer func() {
		if err := os.RemoveAll(COPY_TEST_DIR); err != nil {
			panic(err)
		}
	}()

	if err := exec.Command("go", "build", "-o", fmt.Sprintf("%s/editor", COPY_TEST_DIR), "github.com/jobtalk/pnzr/test/dummy_editor").Run(); err != nil {
		panic(err)
	}

	crypter := cryptex.New(rsa.New(privateKey, publicKey))
	testEditor := &Editor{
		crypter,
		fmt.Sprintf("%s/editor", COPY_TEST_DIR),
	}

	infos, err := ioutil.ReadDir(ORIGIN_TEST_DIR)
	if err != nil {
		panic(err)
	}

	tests := map[string]struct {
		input     string
		want      string
		isEncrypt bool
		err       bool
	}{
		"0.json": {
			`{"foo":"bar", "array": ["foo", "bar", "baz"]}`,
			`{"foo":"bar", "array": ["foo", "bar", "baz"]}`,
			true,
			false,
		},
		"1.json": {
			`{"foo":"bar"}`,
			`{"array": ["foo", "bar", "baz"]}`,
			false,
			true,
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

		if test, ok := tests[info.Name()]; ok {
			if test.isEncrypt {
				err := (&encrypter{crypter}).encrypt(COPY_TEST_DIR + "/" + info.Name())
				if err != nil {
					panic(err)
				}
			}
		}
	}

	for key, test := range tests {
		done := make(chan bool)
		go func(chan bool) {
			err := testEditor.Edit(COPY_TEST_DIR + "/" + key)
			if !test.err && err != nil {
				done <- true
				t.Fatalf("should not be error but: %v", err)
			}
			if test.err && err == nil {
				done <- true
				t.Fatalf("should be error")
			}
			if test.err {
				done <- true
				return
			}

			if test.isEncrypt {
				if err := (&encrypter{crypter}).decrypt(COPY_TEST_DIR + "/" + key); err != nil {
					done <- true
					panic(err)
				}
			}
			editedBin, err := ioutil.ReadFile(COPY_TEST_DIR + "/" + key)
			if err != nil {
				done <- true
				panic(err)
			}

			if !compaireJSON([]byte(test.want), editedBin) {
				done <- true
				t.Log(test.want)
				t.Log(string(editedBin))
				t.Fatalf("not match")
			}

			done <- true
		}(done)

		go func() {
			defer func() {
				done <- true
				err := recover()
				if err != nil {
					panic(err)
				}
			}()
			client := http.Client{}
			if err := ping(); err != nil {
				panic(err)
			}

			req, err := http.NewRequest("POST", "http://localhost:8080/e", strings.NewReader(test.input))
			if err != nil {
				panic(err)
			}
			if _, err := client.Do(req); err != nil {
				panic(err)
			}

			if _, err := client.Get("http://localhost:8080/wq"); err != nil {
				panic(err)
			}
		}()

		<-done
	}
}

func ping() error {
	timeup := make(chan bool)
	go func() {
		time.Sleep(5 * time.Second)
		timeup <- true
	}()

	for {
		select {
		case <-timeup:
			return fmt.Errorf("timeout")
		default:
			_, err := (&http.Client{}).Get("http://localhost:8080")
			if err == nil {
				return nil
			}
		}
	}
}
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
