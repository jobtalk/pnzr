package deploy

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/jobtalk/pnzr/vars"
	"github.com/joho/godotenv"
	"os"
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

func TestDeployCommand_MergeUseDotFile(t *testing.T) {
	testDataDir := vars.TEST_DATA_DIR_ROOT + "/subcmd/deploy/mergeParams/mergeUseDotFile"
	tests := []struct {
		args []string
		want *params
	}{
		{
			[]string{},
			&params{
				kmsKeyID:     aws.String(""),
				file:         aws.String(""),
				profile:      aws.String("default"),
				varsPath:     aws.String(""),
				overrideTag:  aws.String("latest"),
				region:       aws.String(""),
				awsAccessKey: aws.String(""),
				awsSecretKey: aws.String(""),
			},
		},
		{
			[]string{
				"-key_id", "hoge",
			},
			&params{
				kmsKeyID:     aws.String("hoge"),
				file:         aws.String(""),
				profile:      aws.String("default"),
				varsPath:     aws.String(""),
				overrideTag:  aws.String("latest"),
				region:       aws.String(""),
				awsAccessKey: aws.String(""),
				awsSecretKey: aws.String(""),
			},
		},
		{
			[]string{
				"-key_id", "primary-id",
			},
			&params{
				kmsKeyID:     aws.String("primary-id"),
				file:         aws.String(""),
				profile:      aws.String("default"),
				varsPath:     aws.String(""),
				overrideTag:  aws.String("latest"),
				region:       aws.String(""),
				awsAccessKey: aws.String(""),
				awsSecretKey: aws.String(""),
			},
		},
		{
			[]string{},
			&params{
				kmsKeyID:     aws.String("secondary-id"),
				file:         aws.String(""),
				profile:      aws.String("default"),
				varsPath:     aws.String(""),
				overrideTag:  aws.String("latest"),
				region:       aws.String(""),
				awsAccessKey: aws.String(""),
				awsSecretKey: aws.String(""),
			},
		},
		{
			[]string{
				"-profile", "primary-profile",
			},
			&params{
				kmsKeyID:     aws.String(""),
				file:         aws.String(""),
				profile:      aws.String("primary-profile"),
				varsPath:     aws.String(""),
				overrideTag:  aws.String("latest"),
				region:       aws.String(""),
				awsAccessKey: aws.String(""),
				awsSecretKey: aws.String(""),
			},
		},
		{
			[]string{
				"-profile", "primary-profile",
			},
			&params{
				kmsKeyID:     aws.String(""),
				file:         aws.String(""),
				profile:      aws.String("primary-profile"),
				varsPath:     aws.String(""),
				overrideTag:  aws.String("latest"),
				region:       aws.String(""),
				awsAccessKey: aws.String(""),
				awsSecretKey: aws.String(""),
			},
		},
		{
			[]string{},
			&params{
				kmsKeyID:     aws.String(""),
				file:         aws.String(""),
				profile:      aws.String("secondary-profile"),
				varsPath:     aws.String(""),
				overrideTag:  aws.String("latest"),
				region:       aws.String(""),
				awsAccessKey: aws.String(""),
				awsSecretKey: aws.String(""),
			},
		},
	}

	for i, test := range tests {
		func(test struct {
			args []string
			want *params
		}) {
			dotEnvFileName := fmt.Sprintf("%s/%d.env", testDataDir, i)
			envMap, err := godotenv.Read(dotEnvFileName)
			if err != nil {
				panic(err)
			}
			if err := godotenv.Load(dotEnvFileName); err != nil {
				panic(err)
			}
			defer func(m map[string]string) {
				for k := range m {
					os.Unsetenv(k)
				}
			}(envMap)

			deployCmd := &DeployCommand{}
			deployCmd.parseArgs(test.args)
			deployCmd.parseEnv()
			deployCmd.mergeParams()

			got := deployCmd.mergedParams

			if !compareParam(got, test.want) {
				t.Fatalf("want %q, but %q:", test.want, got)
			}

		}(test)
	}
}
