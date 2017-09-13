package deploy

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"testing"
)

var (
	TEST_DIR = os.Getenv("GOPATH") + "/src/github.com/jobtalk/pnzr/test"
)

func init() {
}

func TestCompaireStringPointer(t *testing.T) {
	tests := []struct {
		inS1 *string
		inS2 *string
		want bool
	}{
		{
			nil,
			nil,
			true,
		},
		{
			aws.String("hoge"),
			aws.String("hoge"),
			true,
		},
		{
			nil,
			aws.String("hoge"),
			false,
		},
		{
			aws.String("hoge"),
			nil,
			false,
		},
		{
			aws.String("hgoe"),
			aws.String("huga"),
			false,
		},
		{
			aws.String(""),
			nil,
			true,
		},
		{
			nil,
			aws.String(""),
			true,
		},
	}

	for _, test := range tests {
		got := compaireStringPointer(test.inS1, test.inS2)

		if got != test.want {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}
}

func compaireStringPointer(s1, s2 *string) bool {
	if s1 == nil && s2 != nil {
		return *s2 == ""
	} else if s1 != nil && s2 == nil {
		return *s1 == ""
	} else if s1 == nil && s2 == nil {
		return true
	}
	return *s1 == *s2
}

func compaireParam(p1, p2 *params) bool {
	if p1 == nil && p2 != nil {
		return false
	} else if p1 != nil && p2 == nil {
		return false
	}
	if !compaireStringPointer(p1.kmsKeyID, p2.kmsKeyID) {
		fmt.Println("kms key is not match")
		return false
	}
	if !compaireStringPointer(p1.file, p2.file) {
		fmt.Println("file name is not match")
		return false
	}
	if !compaireStringPointer(p1.profile, p2.profile) {
		fmt.Println("profile name is not match")
		return false
	}
	if !compaireStringPointer(p1.region, p2.region) {
		fmt.Println("region name is not match")
		return false
	}
	if !compaireStringPointer(p1.varsPath, p2.varsPath) {
		fmt.Println("vars path name is not match")
		return false
	}
	if !compaireStringPointer(p1.overrideTag, p2.overrideTag) {
		fmt.Println("tag name is not match")
		return false
	}
	if !compaireStringPointer(p1.awsAccessKey, p2.awsAccessKey) {
		fmt.Println("access key is not match")
		return false
	}
	if !compaireStringPointer(p1.awsSecretKey, p2.awsSecretKey) {
		fmt.Println("secret key is not match")
		return false
	}
	return true
}
