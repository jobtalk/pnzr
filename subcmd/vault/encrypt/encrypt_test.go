package encrypt

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

var (
	TEST_DIR = os.Getenv("GOPATH") + "/src/github.com/jobtalk/pnzr/test"
)

type ApiTest struct{}

func (a *ApiTest) Encrypt(d []byte) ([]byte, error) {
	for i, v := range d {
		d[i] = v + 1
	}

	return d, nil
}

func (a *ApiTest) getReadFile(fileName string) []byte {
	bin, err := ioutil.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	got, err := a.Encrypt(bin)
	if err != nil {
		panic(err)
	}
	return got
}

func TestEncrypt(t *testing.T) {
	var a *ApiTest
	tests := []struct {
		in   string
		want []byte
		err  bool
	}{
		{
			in:   TEST_DIR + "/embedde-test-base.json",
			want: a.getReadFile(TEST_DIR + "/embedde-test-base.json"),
			err:  false,
		},
		{
			in:   TEST_DIR + "/embedde-test-val.json",
			want: a.getReadFile(TEST_DIR + "/embedde-test-val.json"),
			err:  false,
		},
	}

	for _, test := range tests {
		bin, err := ioutil.ReadFile(test.in)
		if err != nil && !test.err {
			panic(err)
		}

		got, err := a.Encrypt(bin)
		if err != nil && !test.err {
			t.Fatalf("should not be error for %v but: %v", test.in, err)
		}

		if test.err && err == nil {
			t.Fatalf("should be error for %v but not:", test.in)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("want %q, but %q", test.want, got)
		}
	}
}
