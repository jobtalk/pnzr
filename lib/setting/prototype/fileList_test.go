package prototype

import (
	"fmt"
	"github.com/jobtalk/pnzr/vars"
	"reflect"
	"testing"
)

func TestFileList(t *testing.T) {
	testDataDir := vars.TEST_DATA_DIR_ROOT + "/subcmd/deploy/testFileList"

	tests := []struct {
		want []string
		err  bool
	}{
		{
			want: func() []string {
				ret := []string{}
				for i := 0; i < 10; i++ {
					ret = append(ret, fmt.Sprintf("%d.json", i))
				}
				return ret
			}(),
			err: false,
		},
		{
			want: []string{},
			err:  false,
		},
		{
			want: func() []string {
				ret := []string{}
				for i := 0; i < 10; i++ {
					ret = append(ret, fmt.Sprintf("%d.json", i))
				}
				for i := 0; i <= int(byte('z')-byte('a')); i++ {
					ret = append(ret, fmt.Sprintf("%s.json", string([]rune{rune(byte('a') + byte(i))})))
				}
				return ret
			}(),
			err: false,
		},
		{
			want: func() []string {
				ret := []string{}
				for i := 0; i < 10; i++ {
					ret = append(ret, fmt.Sprintf("%d.json", i))
				}
				return ret
			}(),
			err: false,
		},
	}

	for i, test := range tests {
		got, err := fileList(fmt.Sprintf("%s/%d", testDataDir, i))
		if !test.err && err != nil {
			t.Fatalf("should not be error for %v but: %v, dir: %v", i, err, fmt.Sprintf("%s/%d", testDataDir, i))
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %v but not:", i)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Fatalf("want: %q, but: %q", test.want, got)
		}
	}
}
