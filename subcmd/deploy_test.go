package subcmd

import (
	"os"
	"testing"
)

var (
	TEST_DIR = os.Getenv("GOPATH") + "/src/github.com/ieee0824/thor/test"
)

func TestReadExternalVariablesFromFile(t *testing.T) {
	var EXTERNSL_DIR = TEST_DIR + "/readExternalVariablesFromFiles"
	result, err := readExternalVariablesFromFile(EXTERNSL_DIR)
	if err != nil {
		t.Error(err)
	} else if len(result) != 2 {
		t.Errorf("The number of elements is invalid. Originally it should be %v, but it is actually %v.", 2, len(result))
	}
}

func TestReadExternalVariables(t *testing.T) {
	var EXTERNSL_DIRS = []string{
		"readExternalVariablesFiles/dir1",
		"readExternalVariablesFiles/dir2",
		"readExternalVariablesFiles/dir3",
	}

	result, err := readExternalVariables(EXTERNSL_DIRS...)
	if err != nil {
		t.Error(err)
	} else if len(result) != 2 {
		t.Errorf("The number of elements is invalid. Originally it should be %v, but it is actually %v.", 2, len(result))
	}
}
