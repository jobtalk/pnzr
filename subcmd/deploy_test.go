package subcmd

import (
	"os"
	"testing"
)

var (
	TEST_DIR     = os.Getenv("GOPATH") + "/src/github.com/ieee0824/thor/test"
	EXTERNSL_DIR = TEST_DIR + "/readExternalVariablesFromFiles"
)

func TestReadExternalVariablesFromFile(t *testing.T) {
	result, err := readExternalVariablesFromFile(EXTERNSL_DIR)
	if err != nil {
		t.Error(err)
	} else if len(result) != 2 {
		t.Errorf("The number of elements is invalid. Originally it should be %v, but it is actually %v.", 2, len(result))
	}
}
