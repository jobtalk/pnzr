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
	tests := []struct {
		argsParams *params
		envParams  *params
		want       *params
	}{
		{
			&params{},
			&params{},
			&params{},
		},
		{
			&params{},
			nil,
			&params{},
		},
		{
			nil,
			&params{},
			&params{},
		},
		{
			nil,
			nil,
			nil,
		},
		{
			&params{
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
			},
			nil,
			&params{
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
			},
		},
		{
			&params{
				aws.String("hoge"),
				aws.String("huga"),
				aws.String("foo"),
				aws.String("bar"),
				aws.String("baz"),
				aws.String("fizz"),
				aws.String("bazz"),
				aws.String("fizzbazz"),
			},
			nil,
			&params{
				aws.String("hoge"),
				aws.String("huga"),
				aws.String("foo"),
				aws.String("bar"),
				aws.String("baz"),
				aws.String("fizz"),
				aws.String("bazz"),
				aws.String("fizzbazz"),
			},
		},
		{
			nil,
			&params{
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
			},
			&params{
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
			},
		},
		{
			nil,
			&params{
				aws.String("hoge"),
				aws.String("huga"),
				aws.String("foo"),
				aws.String("bar"),
				aws.String("baz"),
				aws.String("fizz"),
				aws.String("bazz"),
				aws.String("fizzbazz"),
			},
			&params{
				aws.String("hoge"),
				aws.String("huga"),
				aws.String("foo"),
				aws.String("bar"),
				aws.String("baz"),
				aws.String("fizz"),
				aws.String("bazz"),
				aws.String("fizzbazz"),
			},
		},
		{
			&params{
				aws.String("hoge"),
				aws.String("huga"),
				aws.String("foo"),
				aws.String("bar"),
				aws.String("baz"),
				aws.String("fizz"),
				aws.String("bazz"),
				aws.String("fizzbazz"),
			},
			&params{
				aws.String("1"),
				aws.String("2"),
				aws.String("3"),
				aws.String("4"),
				aws.String("5"),
				aws.String("6"),
				aws.String("7"),
				aws.String("8"),
			},
			&params{
				aws.String("hoge"),
				aws.String("huga"),
				aws.String("foo"),
				aws.String("bar"),
				aws.String("baz"),
				aws.String("fizz"),
				aws.String("bazz"),
				aws.String("fizzbazz"),
			},
		},
		{
			&params{
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
			},
			&params{
				aws.String("1"),
				aws.String("2"),
				aws.String("3"),
				aws.String("4"),
				aws.String("5"),
				aws.String("6"),
				aws.String("7"),
				aws.String("8"),
			},
			&params{
				aws.String("1"),
				aws.String("2"),
				aws.String("3"),
				aws.String("4"),
				aws.String("5"),
				aws.String("6"),
				aws.String("7"),
				aws.String("8"),
			},
		},
		{
			&params{
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
			},
			&params{},
			&params{
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
				aws.String(""),
			},
		},
	}

	for _, test := range tests {
		deployCmd := &DeployCommand{
			paramsFromArgs: test.argsParams,
			paramsFromEnvs: test.envParams,
		}

		deployCmd.mergeParams()

		got := deployCmd.mergedParams

		if test.argsParams == nil && test.envParams == nil && got != nil {
			t.Fatalf("error prams is not nil")
		}
		if got != nil && !compareParam(got, test.want) {
			t.Fatalf("error prams is not match")
		}
	}
}
