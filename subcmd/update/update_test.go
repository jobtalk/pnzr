package update

import (
	"log"
	"testing"

	"github.com/jobtalk/pnzr/vars"
)

func TestCheckVersion(t *testing.T) {
	vars.VERSION = "v1.2.0"
	tests := []struct {
		in    string
		want1 string
		want2 int
	}{
		{
			"v1.2.0",
			"this version is latest",
			0,
		},
		{
			"",
			"can not get latest versiont",
			255,
		},
		{
			"hoge",
			"",
			0,
		},
	}

	for _, test := range tests {
		gotInt, gotS := checkVersion(test.in)
		if gotInt != test.want2 || gotS != test.want1 {
			log.Println(vars.VERSION)

			t.Fatalf("want %s, %d, but %s, %d:", test.want1, test.want2, gotS, gotInt)
		}
	}
}
