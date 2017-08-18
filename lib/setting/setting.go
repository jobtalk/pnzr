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
	globalSession *session.Session
	option = struct {
		region *string
		profile *string
	}{}
)

func roundArgs(args []string) (roundedArgs []string, region string, profile string) {
	profile = getenv.String("AWS_PROFILE_NAME", "default")
	region = getenv.String("AWS_DEFAULT_REGION")

	for i := 0; i < len(args); i++ {
		if args[i] == "-region" {
			region = args[i+1]
			i++
		} else if args[i] == "-profile" {
			profile = args[i+1]
			i++
		} else if strings.HasPrefix(args[i], "-region=") {
			region = strings.TrimPrefix(args[i], "-region=")
		} else if strings.HasPrefix(args[i], "-profile=") {
			profile = strings.TrimPrefix(args[i], "-profile=")
		} else {
			roundedArgs = append(roundedArgs, args[i])
		}
	}

	return
}

func Initial(args []string) ([]string, error) {
	args, region, profile := roundArgs(args)
	option.region = &region
	option.profile = &profile

	globalSession = session.Must(session.NewSessionWithOptions(session.Options{
		AssumeRoleTokenProvider: stscreds.StdinTokenProvider,
		SharedConfigState:       session.SharedConfigEnable,
		Profile:                 *option.profile,
		Config: aws.Config{
			Region: option.region,
		},
	}))

	return args, nil
}

func Region()*string{
	return option.region
}

func Profile()*string{
	return option.profile
}

func GetSession() *session.Session {
	return globalSession
}