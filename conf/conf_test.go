package conf

import (
	"io/ioutil"
	"os"
	"testing"
)

var (
	TEST_DIR = os.Getenv("GOPATH") + "/src/github.com/jobtalk/thor/test"
)

func TestIsJSON(t *testing.T) {
	jsonStr := `{"hoge":"huga", "flag": true, "num": 0}`
	nonJSON := "This is a pen."

	if !isJSON(jsonStr) {
		t.Errorf("The expected value is %v, but actually it is %v.", true, !isJSON(jsonStr))
	}

	if isJSON(nonJSON) {
		t.Errorf("The expected value is %v, but actually it is %v.", false, isJSON(nonJSON))
	}
}

func TestEmbedde(t *testing.T) {
	baseJSON, err := ioutil.ReadFile(TEST_DIR + "/" + "embedde-test-base.json")
	if err != nil {
		t.Errorf("base test data can not read, error is : %v", err)
	}
	valJSON, err := ioutil.ReadFile(TEST_DIR + "/" + "embedde-test-val.json")
	if err != nil {
		t.Errorf("val test data can not read, error is : %v", err)
	}

	conf, err := Embedde(string(baseJSON), string(valJSON))
	if err != nil {
		t.Error(err)
	}
	if conf != `{"hoge":{"hoge":"huga"}}` {
		t.Errorf("The expected value is `%v`, but actually it is `%v`.", `{"hoge":{"hoge":"huga"}}`, conf)
	}
}
