package decrypt

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

type ApiTest struct{}

var (
	TEST_DIR = os.Getenv("GOPATH") + "/src/github.com/jobtalk/pnzr/test"
)

func (a *ApiTest) Encrypt(d []byte) ([]byte, error) {
	for i, v := range d {
		d[i] = v + 1
	}

	return d, nil
}

func (a *ApiTest) Decrypt(d []byte) ([]byte, error) {
	for i, v := range d {
		d[i] = v - 1
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

func TestDecrypt(t *testing.T) {
	var a *ApiTest
	tests := []struct {
		in   []byte
		want string
		err  bool
	}{
		{
			in:   a.getReadFile(TEST_DIR + "/embedde-test-base.json"),
			want: TEST_DIR + "/embedde-test-base.json",
			err:  false,
		},
		{
			in:   a.getReadFile(TEST_DIR + "/embedde-test-val.json"),
			want: TEST_DIR + "/embedde-test-val.json",
			err:  false,
		},
	}
	for _, test := range tests {
		got, err := a.Decrypt(test.in)
		if err != nil && !test.err {
			t.Fatalf("should not be error for %v but: %v", test.in, err)
		}

		if test.err && err == nil {
			t.Fatalf("should be error for %v but not:", test.in)
		}

		want, err := ioutil.ReadFile(test.want)
		if err != nil {
			panic(err)
		}

		if !reflect.DeepEqual(got, want) {
			t.Fatalf("want %q, but %q:", test.want, got)
		}

	}

}
