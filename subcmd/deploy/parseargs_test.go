package deploy

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
)

func TestDeployCommand_parseArgs(t *testing.T) {
	tests := []struct {
		input      []string
		want       *params
		help       bool
		expectsErr bool
	}{
		{
			input:      []string{},
			want:       &params{},
			help:       false,
			expectsErr: false,
		},
		{
			input:      []string{"-h"},
			want:       &params{},
			help:       true,
			expectsErr: false,
		},
		{
			input: []string{
				"-key_id",
				"hoge",
			},
			want: &params{
				kmsKeyID: aws.String("hoge"),
			},
			help:       false,
			expectsErr: false,
		},
		{
			input: []string{
				"-key_id",
			},
			help:       false,
			expectsErr: true,
		},
		{
			input: []string{
				"-f",
				"hoge",
			},
			want: &params{
				file: aws.String("hoge"),
			},
			help:       false,
			expectsErr: false,
		},
		{
			input: []string{
				"-f",
			},
			help:       false,
			expectsErr: true,
		},
		{
			input: []string{
				"-file",
				"huga",
			},
			want: &params{
				file: aws.String("huga"),
			},
			help:       false,
			expectsErr: false,
		},
		{
			input: []string{
				"-file",
			},
			help:       false,
			expectsErr: true,
		},
		{
			input: []string{
				"foo",
			},
			want: &params{
				file: aws.String("foo"),
			},
			help:       false,
			expectsErr: false,
		},
		{
			input: []string{
				"-f", "hoge",
				"-file", "huga",
			},
			want: &params{
				file: aws.String("huga"),
			},
			help:       false,
			expectsErr: false,
		},
		{
			input: []string{
				"-f", "hoge",
				"-file",
			},
			help:       false,
			expectsErr: true,
		},
		{
			input: []string{
				"-f", "hoge",
				"-file", "huga",
				"foo",
			},
			want: &params{
				file: aws.String("huga"),
			},
			help:       false,
			expectsErr: false,
		},
		{
			input: []string{
				"-profile", "ap-northeast-1",
			},
			want: &params{
				profile: aws.String("ap-northeast-1"),
			},
			help:       false,
			expectsErr: false,
		},
		{
			input: []string{
				"-vars_path", "hoge",
			},
			want: &params{
				varsPath: aws.String("hoge"),
			},
			help:       false,
			expectsErr: false,
		},
		{
			input: []string{
				"-t", "hoge",
			},
			want: &params{
				overrideTag: aws.String("hoge"),
			},
			help:       false,
			expectsErr: false,
		},
		{
			input: []string{
				"-aws-access-key-id", "hoge",
			},
			want: &params{
				awsAccessKey: aws.String("hoge"),
			},
			help:       false,
			expectsErr: false,
		},
		{
			input: []string{
				"-aws-secret-key-id", "hoge",
			},
			want: &params{
				awsSecretKey: aws.String("hoge"),
			},
			help:       false,
			expectsErr: false,
		},
		{
			input: []string{
				"-aaaaaaaa",
			},
			help:       false,
			expectsErr: true,
		},
		{
			input: []string{
				"-key_id", "some_key",
				"-f", "some_file",
				"-profile", "some_profile",
				"-region", "some_region",
				"-vars_path", "some_vars_path",
				"-t", "some_tag",
				"-aws-access-key-id", "access_key",
				"-aws-secret-key-id", "secret_key",
			},
			want: &params{
				kmsKeyID:     aws.String("some_key"),
				file:         aws.String("some_file"),
				profile:      aws.String("some_profile"),
				region:       aws.String("some_region"),
				varsPath:     aws.String("some_vars_path"),
				overrideTag:  aws.String("some_tag"),
				awsAccessKey: aws.String("access_key"),
				awsSecretKey: aws.String("secret_key"),
			},
			help:       false,
			expectsErr: false,
		},
	}

	for _, test := range tests {
		func(test struct {
			input      []string
			want       *params
			help       bool
			expectsErr bool
		}) {
			defer func() {
				err := recover()
				if !test.expectsErr && err != nil {
					t.Fatalf("should not be error for %v but: %v", test.input, err)
				}
				if test.expectsErr && err == nil {
					t.Fatalf("should be error for %v but not:", test.input)
				}
			}()
			deployCmd := &DeployCommand{}
			help := deployCmd.parseArgs(test.input)
			if !test.help && help != "" {
				t.Fatalf("should not be help for %v but: %v", test.input, help)
			}
			if test.help && help == "" {
				t.Fatalf("should be help for %v but not:", test.input)
			}
			if help != "" {
				return
			}
			if !compareParam(deployCmd.paramsFromArgs, test.want) {
				t.Fatalf("error prams is not match")
			}

		}(test)
	}
}
