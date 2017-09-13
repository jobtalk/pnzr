package deploy

import (
	"github.com/aws/aws-sdk-go/aws"
	"testing"
)

func TestStringIsEmpty(t *testing.T) {
	tests := []struct {
		input *string
		want  bool
	}{
		{
			nil,
			true,
		},
		{
			aws.String(""),
			true,
		},
		{
			aws.String("hoge"),
			false,
		},
	}

	for _, test := range tests {
		got := stringIsEmpty(test.input)
		if test.want != got {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}
}

func TestDeployCommand_MergeParams(t *testing.T) {

}
