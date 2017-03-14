package subcmd

import (
	"testing"
)

func TestGetFullNameParam(t *testing.T) {
	testArgs := []string{
		"--file=test.txt",
		"--hoge=hoge.txt",
		"--hoge=huga.txt",
		"--foo",
	}
	vals, err := getFullNameParam(testArgs, "--file")
	if len(vals) != 1 {
		t.Error("Variable number of elements. Number of elements should be 1.: %v", len(vals))
	} else if err != nil {
		t.Error(err)
	} else if vals[0] == nil {
		t.Errorf("args val is %v", vals[0])
	} else if *vals[0] != "test.txt" {
		t.Errorf("args val is %v", *vals[0])
	}

	vals, err = getFullNameParam(testArgs, "--hoge")
	if len(vals) != 2 {
		t.Error("Variable number of elements. Number of elements should be 2.: %v", len(vals))
	} else if err != nil {
		t.Error(err)
	} else if vals[0] == nil {
		t.Errorf("args val is %v", vals[0])
	} else if *vals[0] != "hoge.txt" {
		t.Errorf("args val is %v", *vals[0])
	} else if vals[1] == nil {
		t.Errorf("args val is %v", vals[1])
	} else if *vals[1] != "huga.txt" {
		t.Errorf("args val is %v", *vals[1])
	}

	vals, err = getFullNameParam(testArgs, "--foo")
	if err != nil {
		t.Error(err)
	} else if *vals[0] != "true" {
		t.Error("Invalid parameter: %v", *vals[0])
	}
}
