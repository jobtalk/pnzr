package v1

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/ieee0824/cryptex"
	"github.com/ieee0824/cryptex/rsa"
	"github.com/jobtalk/pnzr/vars"
)

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

type editTester struct {
	ORIGIN_TEST_DIR string
	COPY_TEST_DIR   string
	KEY_DIR         string
	privateKey      []byte
	publicKey       []byte
	cryptex         *cryptex.Cryptex
	editor          *Editor
	encrypter       *encrypter
}

func newEditTester() (*editTester, error) {
	ret := &editTester{}
	ret.ORIGIN_TEST_DIR = vars.TEST_DATA_DIR_ROOT + "/subcmd/vault/edit/v1"
	ret.COPY_TEST_DIR = ret.ORIGIN_TEST_DIR + "/copy"
	ret.KEY_DIR = vars.TEST_DATA_DIR_ROOT + "/key"

	if err := ret.generateRsaKey(); err != nil {
		return nil, err
	}
	ret.generateCryptex()
	ret.encrypter = &encrypter{ret.cryptex}

	if err := ret.generateTestDir(); err != nil {
		return nil, err
	}

	if err := ret.generateEditor(); err != nil {
		return nil, err
	}

	return ret, nil
}

func (t *editTester) generateRsaKey() error {
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

func (t *editTester) generateCryptex() {
	t.cryptex = cryptex.New(rsa.New(t.privateKey, t.publicKey))
}

func (t *editTester) generateTestDir() error {
	if err := exec.Command("mkdir", "-p", t.COPY_TEST_DIR).Run(); err != nil {
		return err
	}
	return nil
}

func (t *editTester) fin() error {
	if err := os.RemoveAll(t.COPY_TEST_DIR); err != nil {
		return err
	}
	return nil
}

func (t *editTester) generateEditor() error {
	if err := exec.Command("go", "build", "-o", fmt.Sprintf("%s/editor", t.COPY_TEST_DIR), "github.com/jobtalk/pnzr/test/dummy_editor").Run(); err != nil {
		return err
	}
	t.editor = &Editor{
		t.cryptex,
		fmt.Sprintf("%s/editor", t.COPY_TEST_DIR),
	}
	return nil
}

func (t *editTester) generateEncryptedFiles(tests map[string]testInput) error {
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

		if test, ok := tests[info.Name()]; ok {
			if test.encryptsTestFileAtFirst {
				err := t.encrypter.encrypt(t.COPY_TEST_DIR + "/" + info.Name())
				if err != nil {
					panic(err)
				}
			}
		}
	}
	return nil
}

func (tester *editTester) execEdit(t *testing.T, key string, test testInput, done chan<- bool, errCh chan<- error) {
	defer func() {
		done <- true
		close(done)
	}()
	err := tester.editor.Edit(tester.COPY_TEST_DIR + "/" + key)
	if !test.expectsErr && err != nil {
		t.Fatalf("should not be error but: %v", err)
		return
	}
	if test.expectsErr && err == nil {
		t.Fatalf("should be error")
	}

	if test.expectsErr {
		return
	}

	if test.encryptsTestFileAtFirst {
		if err := tester.encrypter.decrypt(tester.COPY_TEST_DIR + "/" + key); err != nil {
			errCh <- err
			return
		}
	}
	editedBin, err := ioutil.ReadFile(tester.COPY_TEST_DIR + "/" + key)
	if err != nil {
		errCh <- err
		return
	}
	if !compareJSON([]byte(test.want), editedBin) {
		t.Log(test.want)
		t.Log(string(editedBin))
		t.Fatalf("not match")
	}
	return
}

func sendEditorCommand(errCh chan error, test testInput) {
	client := http.Client{}

	if err := ping(); err != nil {
		errCh <- err
		return
	}

	req, err := http.NewRequest("POST", "http://localhost:8080/e", strings.NewReader(test.input))
	if err != nil {
		errCh <- err
		return
	}

	if _, err := client.Do(req); err != nil {
		errCh <- err
		return
	}
	if _, err := client.Get("http://localhost:8080/wq"); err != nil {
		errCh <- err
		return
	}
}

func (tester *editTester) run(t *testing.T, tests map[string]testInput) error {
	for key, test := range tests {
		done := make(chan bool)
		errCh := make(chan error)
		go tester.execEdit(t, key, test, done, errCh)
		go sendEditorCommand(errCh, test)

		select {
		case err := <-errCh:
			return err
		case <-done:
		}
	}

	return nil
}

// testInput struct defines test data.
//
// input field is the input value of the test.
// This test corresponds to the contents edited by the Editor.
//
// want field is the expected value after the test.
//
// encryptsTestFileAtFirst field is a flag
// that controls whether to encrypt at the start of the test.
//
// expectsErr field is a flag
// that controls whether or not to expect errors.
type testInput struct {
	input                   string
	want                    string
	encryptsTestFileAtFirst bool
	expectsErr              bool
}

func TestEditor_Edit1(t *testing.T) {
	tester, err := newEditTester()
	if err != nil {
		panic(err)
	}

	defer func() {
		err := tester.fin()
		if err != nil {
			panic(err)
		}
	}()
	tests := map[string]testInput{
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

	if err := tester.generateEncryptedFiles(tests); err != nil {
		panic(err)
	}

	if err := tester.run(t, tests); err != nil {
		panic(err)
	}
}

func TestEditor_Edit2(t *testing.T) {
	tester, err := newEditTester()
	if err != nil {
		panic(err)
	}
	defer func() {
		err := tester.fin()
		if err != nil {
			panic(err)
		}
	}()

	tester.editor = &Editor{
		cryptex.New(newChiper(true)),
		fmt.Sprintf("%s/editor", tester.COPY_TEST_DIR),
	}

	tests := []struct {
		fileName string
	}{
		{"hoge"},
		{tester.ORIGIN_TEST_DIR + "/0.json"},
	}

	for _, test := range tests {
		err := tester.editor.Edit(test.fileName)
		if err == nil {
			t.Fatalf("should be error for %v but not:", test.fileName)
		}
	}
}
