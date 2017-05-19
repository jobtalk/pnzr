package lib

import (
	_ "fmt"
	"testing"

	ecssdk "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

func TestListServices(t *testing.T) {
	ecs, _ := mockECS()
	res, _ := ecs.ListServices(&ecssdk.ListServicesInput{})
	t.Log(res)
	for i, arn := range res.ServiceArns {
		if *testServiceArns[i] != *arn {
			t.Fatalf("Expected: %s, Got: %s", *testServiceArns[i], *arn)
		}
	}
}

func TestUpsertService(t *testing.T) {
	// 存在するサービスは Update が呼ばれる
	{
		svcName := "some-service-a"
		svc := &ecssdk.CreateServiceInput{ServiceName: &svcName}
		ecs, fnArgs := mockECS()
		ecs.UpsertService(svc)
		if *fnArgs.UpdateServiceInput.Service != svcName {
			t.Log(*fnArgs.UpdateServiceInput.Service)
			t.Fatalf("UpdateService should be called with %s", svcName)
		}
		if fnArgs.CreateServiceInput != nil {
			t.Fatalf("CreateService should not be called")
		}
	}

	// 存在しないサービスは Create が呼ばれる
	{
		svcName := "some-service-d"
		svc := &ecssdk.CreateServiceInput{ServiceName: &svcName}
		ecs, fnArgs := mockECS()
		ecs.UpsertService(svc)
		if *fnArgs.CreateServiceInput.ServiceName != svcName {
			t.Log(*fnArgs.CreateServiceInput.ServiceName)
			t.Fatalf("CreateService should be called with %s", svcName)
		}
		if fnArgs.UpdateServiceInput != nil {
			t.Fatalf("UpdateService should not be called")
		}
	}
}

type mockedFnArgs struct {
	CreateServiceInput *ecssdk.CreateServiceInput
	UpdateServiceInput *ecssdk.UpdateServiceInput
}

func mockECS() (*ECS, *mockedFnArgs) {
	fnArgs := mockedFnArgs{CreateServiceInput: nil, UpdateServiceInput: nil}
	svc := &mockedECS{
		listServicesOutput: &ecssdk.ListServicesOutput{ServiceArns: testServiceArns},
		createService: func(in *ecssdk.CreateServiceInput) (*ecssdk.CreateServiceOutput, error) {
			fnArgs = mockedFnArgs{CreateServiceInput: in, UpdateServiceInput: fnArgs.UpdateServiceInput}
			return nil, nil
		},
		updateService: func(in *ecssdk.UpdateServiceInput) (*ecssdk.UpdateServiceOutput, error) {
			fnArgs = mockedFnArgs{CreateServiceInput: fnArgs.CreateServiceInput, UpdateServiceInput: in}
			return nil, nil
		},
	}
	return &ECS{svc: svc}, &fnArgs
}

type mockedECS struct {
	ecsiface.ECSAPI
	listServicesOutput *ecssdk.ListServicesOutput
	createService      func(*ecssdk.CreateServiceInput) (*ecssdk.CreateServiceOutput, error)
	updateService      func(*ecssdk.UpdateServiceInput) (*ecssdk.UpdateServiceOutput, error)
}

func (m mockedECS) CreateService(in *ecssdk.CreateServiceInput) (*ecssdk.CreateServiceOutput, error) {
	m.createService(in)
	return nil, nil
}

func (m mockedECS) UpdateService(in *ecssdk.UpdateServiceInput) (*ecssdk.UpdateServiceOutput, error) {
	m.updateService(in)
	return nil, nil
}

func (m mockedECS) ListServicesPages(in *ecssdk.ListServicesInput, f func(*ecssdk.ListServicesOutput, bool) bool) error {
	f(m.listServicesOutput, true)
	return nil
}

// 既存サービスのモック
var testServiceArns = make([]*string, 3)

func init() {
	s0 := "arn:aws:ecs:xx-someregion-1:01234567890123:service/some-service-a"
	s1 := "arn:aws:ecs:xx-someregion-1:01234567890123:service/some-service-b"
	s2 := "arn:aws:ecs:xx-someregion-1:01234567890123:service/some-service-c"
	testServiceArns[0] = &s0
	testServiceArns[1] = &s1
	testServiceArns[2] = &s2
}
