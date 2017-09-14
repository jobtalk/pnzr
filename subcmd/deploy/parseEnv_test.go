package deploy

import (
	"github.com/aws/aws-sdk-go/aws"
	"os"
	"testing"
)

func eraseEnv() {
	if err := os.Unsetenv("KMS_KEY_ID"); err != nil {
		panic(err)
	}
	if err := os.Unsetenv("AWS_PROFILE_NAME"); err != nil {
		panic(err)
	}
	if err := os.Unsetenv("DOCKER_DEFAULT_DEPLOY_TAG"); err != nil {
		panic(err)
	}
	if err := os.Unsetenv("AWS_REGION"); err != nil {
		panic(err)
	}
	if err := os.Unsetenv("AWS_ACCESS_KEY_ID"); err != nil {
		panic(err)
	}
	if err := os.Unsetenv("AWS_SECRET_ACCESS_KEY"); err != nil {
		panic(err)
	}
}

func TestDeployCommand_ParseEnv(t *testing.T) {
	eraseEnv()
	tests := []struct {
		envs map[string]string
		want *params
	}{
		{
			map[string]string{},
			&params{
				kmsKeyID:     aws.String(""),
				profile:      aws.String("default"),
				overrideTag:  aws.String("latest"),
				region:       aws.String(""),
				awsAccessKey: aws.String(""),
				awsSecretKey: aws.String(""),
			},
		},
		{
			map[string]string{
				"KMS_KEY_ID": "hoge",
			},
			&params{
				kmsKeyID:     aws.String("hoge"),
				profile:      aws.String("default"),
				overrideTag:  aws.String("latest"),
				region:       aws.String(""),
				awsAccessKey: aws.String(""),
				awsSecretKey: aws.String(""),
			},
		},
		{
			map[string]string{
				"AWS_PROFILE_NAME": "pnzr",
			},
			&params{
				kmsKeyID:     aws.String(""),
				profile:      aws.String("pnzr"),
				overrideTag:  aws.String("latest"),
				region:       aws.String(""),
				awsAccessKey: aws.String(""),
				awsSecretKey: aws.String(""),
			},
		},
		{
			map[string]string{
				"DOCKER_DEFAULT_DEPLOY_TAG": "pnzr",
			},
			&params{
				kmsKeyID:     aws.String(""),
				profile:      aws.String("default"),
				overrideTag:  aws.String("pnzr"),
				region:       aws.String(""),
				awsAccessKey: aws.String(""),
				awsSecretKey: aws.String(""),
			},
		},
		{
			map[string]string{
				"AWS_REGION": "ap-northeast-1",
			},
			&params{
				kmsKeyID:     aws.String(""),
				profile:      aws.String("default"),
				overrideTag:  aws.String("latest"),
				region:       aws.String("ap-northeast-1"),
				awsAccessKey: aws.String(""),
				awsSecretKey: aws.String(""),
			},
		},
		{
			map[string]string{
				"AWS_ACCESS_KEY_ID": "foo",
			},
			&params{
				kmsKeyID:     aws.String(""),
				profile:      aws.String("default"),
				overrideTag:  aws.String("latest"),
				region:       aws.String(""),
				awsAccessKey: aws.String("foo"),
				awsSecretKey: aws.String(""),
			},
		},
		{
			map[string]string{
				"AWS_SECRET_ACCESS_KEY": "bar",
			},
			&params{
				kmsKeyID:     aws.String(""),
				profile:      aws.String("default"),
				overrideTag:  aws.String("latest"),
				region:       aws.String(""),
				awsAccessKey: aws.String(""),
				awsSecretKey: aws.String("bar"),
			},
		},
		{
			map[string]string{
				"HOGE": "huga",
			},
			&params{
				kmsKeyID:     aws.String(""),
				profile:      aws.String("default"),
				overrideTag:  aws.String("latest"),
				region:       aws.String(""),
				awsAccessKey: aws.String(""),
				awsSecretKey: aws.String(""),
			},
		},
		{
			map[string]string{
				"KMS_KEY_ID":                "kms",
				"AWS_PROFILE_NAME":          "profile",
				"DOCKER_DEFAULT_DEPLOY_TAG": "tag",
				"AWS_REGION":                "region",
				"AWS_ACCESS_KEY_ID":         "access_key",
				"AWS_SECRET_ACCESS_KEY":     "secret_key",
			},
			&params{
				kmsKeyID:     aws.String("kms"),
				profile:      aws.String("profile"),
				overrideTag:  aws.String("tag"),
				region:       aws.String("region"),
				awsAccessKey: aws.String("access_key"),
				awsSecretKey: aws.String("secret_key"),
			},
		},
	}

	for _, test := range tests {
		func(test struct {
			envs map[string]string
			want *params
		}) {
			defer func(m map[string]string) {
				for key := range m {
					err := os.Unsetenv(key)
					if err != nil {
						panic(err)
					}
				}
			}(test.envs)
			for key, val := range test.envs {
				if err := os.Setenv(key, val); err != nil {
					panic(err)
				}
			}

			deployCmd := &DeployCommand{}
			deployCmd.parseEnv()
			got := deployCmd.paramsFromEnvs

			if !compareParam(got, test.want) {
				t.Fatalf("error prams is not match")
			}
		}(test)
	}
}
