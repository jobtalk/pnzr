package update

import (
	"log"
	"testing"

	"github.com/jobtalk/pnzr/vars"
)

func TestCheckVersion(t *testing.T) {
	vars.VERSION = "v1.2.0"
	tests := []struct {
		inS1    string
		wantS   string
		wantInt int
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
		gotInt, gotS := checkVersion(test.inS1)
		if gotInt != test.wantInt || gotS != test.wantS {
			log.Println(vars.VERSION)

			t.Fatalf("want %s, %d, but %s, %d:", test.wantS, test.wantInt, gotS, gotInt)
		}
	}
}
