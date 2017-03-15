package vault

import (
	"io/ioutil"
	"os"
	"testing"
)

var (
	TEST_DIR = os.Getenv("GOPATH") + "/src/github.com/ieee0824/thor/test"
)

func TestIsSecret(t *testing.T) {
	plainTextPath := TEST_DIR + "/vaultTestFiles/plain.json"
	chipherTextPath := TEST_DIR + "/vaultTestFiles/chipher.json"

	plain, err := ioutil.ReadFile(plainTextPath)
	if err != nil {
		t.Error(err)
	} else if IsSecret(plain) {
		t.Errorf("The result is illegal. I want %v, but it is actually %v.", false, IsSecret(plain))
	}

	chipher, err := ioutil.ReadFile(chipherTextPath)
	if err != nil {
		t.Error(err)
	} else if !IsSecret(chipher) {
		t.Errorf("The result is illegal. I want %v, but it is actually %v.", true, IsSecret(chipher))
	}
}
