package setting

import (
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"strings"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/ieee0824/getenv"
)

type ELB struct {
	LB          *elbv2.CreateLoadBalancerInput
	TargetGroup *elbv2.CreateTargetGroupInput
	Listener    *elbv2.CreateListenerInput
}

type ECS struct {
	Service        *ecs.CreateServiceInput
	TaskDefinition *ecs.RegisterTaskDefinitionInput
}

type Setting struct {
	ELB *ELB
	ECS *ECS
}


var (
	s *session.Session
	option = struct {
		region *string
		profile *string
	}{}
)

func roundFlags(s []string) (o []string, region string, profile string) {
	profile = "default"
	region = "ap-northeast-1"
	for i := 0; i < len(s); i++ {
		if s[i] == "-region" {
			region = s[i+1]
			i++
		} else if s[i] == "-profile" {
			profile = s[i+1]
			i++
		} else if strings.HasPrefix(s[i], "-region=") {
			region = strings.TrimPrefix(s[i], "-region=")
		} else if strings.HasPrefix(s[i], "-profile=") {
			profile = strings.TrimPrefix(s[i], "-profile=")
		} else {
			o = append(o, s[i])
		}
	}

	profile = getenv.String("AWS_PROFILE_NAME", profile)
	region = getenv.String("AWS_REGION", region)
	return
}

func Initial(o []string) ([]string, error) {
	o, region, profile := roundFlags(o)
	option.region = &region
	option.profile = &profile

	s = session.Must(session.NewSessionWithOptions(session.Options{
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
		Profile:                 *option.profile,
		Config: aws.Config{
			Region: option.region,
		},
	}))

	return o, nil
}

func Region()*string{
	return option.region
}

func Profile()*string{
	return option.profile
}

func GetSession() *session.Session {
	return s
}