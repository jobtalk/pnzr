package update

import (
	"fmt"
	"log"
	"testing"

	"github.com/jobtalk/pnzr/vars"
)

func TestCheckVersion(t *testing.T) {
	vars.VERSION = "v1.2.0"
	tests := []struct {
		version        string
		versionMessage string
		exitCode       int
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
			"hoge",
			-1,
		},
	}

	for _, test := range tests {
		gotInt, gotS := checkVersion(test.version)
		if gotInt != test.exitCode || gotS != test.versionMessage {
			log.Println(vars.VERSION)

			t.Fatalf("want %s, %d, but %s, %d:", test.versionMessage, test.exitCode, gotS, gotInt)
		}
	}
}

func TestCheckENV(t *testing.T) {
	tests := []struct {
		os       string
		arch     string
		isUpdate bool
	}{
		{
			"darwin",
			"amd64",
			true,
		},
		{
			"linux",
			"amd64",
			true,
		},
		{
			"hoge",
			"amd64",
			false,
		},
		{
			"darwin",
			"hoge",
			false,
		},
		{
			"linux",
			"hoge",
			false,
		},
		{
			"",
			"",
			false,
		},
	}
	for _, test := range tests {
		e := Env{
			test.os,
			test.arch,
		}
		if ok := e.checkENV(); ok != test.isUpdate {
			t.Errorf("Update checking want to get %v, but %v", test.isUpdate, ok)
		}
	}
}

func TestDetectPlatform(t *testing.T) {
	tests := []struct {
		os       string
		platform string
		err      error
	}{
		{
			"darwin",
			"darwin-amd64",
			nil,
		},
		{
			"linux",
			"linux-amd64",
			nil,
		},
		{
			"hoge",
			"",
			fmt.Errorf("This is not %s", "darwin or linux"),
		},
		{
			"",
			"",
			fmt.Errorf("This is not %s", "darwin or linux"),
		},
	}
	for _, test := range tests {
		e := Env{
			test.os,
			"",
		}
		platform, err := e.detectPlatform()
		if platform != test.platform && err != test.err {
			t.Errorf("want to get %v and %q, but %v and %q", test.platform, test.err, platform, err)
		}
	}
}
