package deploy

import (
	"os"
	"testing"
)

var (
	TEST_DIR = os.Getenv("GOPATH") + "/src/github.com/jobtalk/thor/test"
)

func init() {
}

func TestReadExternalVariablesFromFile(t *testing.T) {
	var EXTERNSL_DIR = TEST_DIR + "/readExternalVariablesFromFiles"
	t.Log("ディレクトリの中身がある時のテスト")
	result, err := readExternalVariablesFromFile(EXTERNSL_DIR)
	if err != nil {
		t.Error(err)
	} else if len(result) != 2 {
		t.Errorf("The number of elements is invalid. Originally it should be %v, but it is actually %v.", 2, len(result))
	} else {
		t.Log("ok")
	}

	t.Log("ディレクトリの中身が無いときのテスト")
	EXTERNSL_DIR = TEST_DIR + "/readExternalVariablesFromFiles/empty"

	os.RemoveAll(EXTERNSL_DIR)
	if err := os.Mkdir(EXTERNSL_DIR, 0777); err != nil {
		t.Error(err)
	}
	result, err = readExternalVariablesFromFile(EXTERNSL_DIR)
	if err != nil {
		t.Error(err)
	} else if len(result) != 0 {
		t.Errorf("The number of elements is invalid. Originally it should be %v, but it is actually %v.", 0, len(result))
	} else {
		t.Log("ok")
	}
	os.RemoveAll(EXTERNSL_DIR)
}

func TestReadExternalVariables(t *testing.T) {
	var EXTERNSL_DIRS = []string{
		TEST_DIR + "/readExternalVariablesFiles/dir1",
		TEST_DIR + "/readExternalVariablesFiles/dir2",
		TEST_DIR + "/readExternalVariablesFiles/dir3",
	}

	os.RemoveAll(TEST_DIR + "/readExternalVariablesFiles/dir3")
	// test用にディレクトリを作る
	if err := os.Mkdir(TEST_DIR+"/readExternalVariablesFiles/dir3", 0777); err != nil {
		t.Error(err)
	}

	result, err := readExternalVariables(EXTERNSL_DIRS...)
	if err != nil {
		t.Error(err)
	} else if len(result) != 2 {
		t.Errorf("The number of elements is invalid. Originally it should be %v, but it is actually %v.", 2, len(result))
	}
	os.RemoveAll(TEST_DIR + "/readExternalVariablesFiles/dir3")
}
