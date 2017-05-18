package lib

import (
	_ "fmt"
	"testing"

	ecssdk "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

func TestListServices(t *testing.T) {
	ecs := newMockedECS()
	res, _ := ecs.ListServices(&ecssdk.ListServicesInput{})
	t.Log(res)
	for i, arn := range res.ServiceArns {
		if *testServiceArns[i] != *arn {
			t.Fatalf("Expected: %s, Got: %s", *testServiceArns[i], *arn)
		}
	}
}

func TestUpsertService(t *testing.T) {
	// 存在するサービスは Update
	existingSvcName := "some-service-a"
	existingSvc := &ecssdk.CreateServiceInput{ServiceName: &existingSvcName}
	{
		ecs, called := updateServiceShouldCalledWith(t, "some-service-a")
		ecs.UpsertService(existingSvc)
		if !*called {
			t.Fatalf("UpdateService should called")
		}
	}
	createServiceShouldNotCalled(t).UpsertService(existingSvc)

	// 存在しないサービスは Create
	newSvcName := "some-service-d"
	newSvc := &ecssdk.CreateServiceInput{ServiceName: &newSvcName}
	{
		ecs, called := createServiceShouldCalledWith(t, "some-service-d")
		ecs.UpsertService(newSvc)
		if !*called {
			t.Fatalf("UpdateService should called")
		}
	}
	updateServiceShouldNotCalled(t).UpsertService(newSvc)
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

func newMockedECS() *ECS {
	svc := &mockedECS{
		listServicesOutput: &ecssdk.ListServicesOutput{ServiceArns: testServiceArns},
		createService:      emptyCreateService,
		updateService:      emptyUpdateService,
	}
	return &ECS{svc: svc}
}

func createServiceShouldCalledWith(t *testing.T, name string) (*ECS, *bool) {
	called := false
	svc := &mockedECS{
		listServicesOutput: &ecssdk.ListServicesOutput{ServiceArns: testServiceArns},
		createService: func(in *ecssdk.CreateServiceInput) (*ecssdk.CreateServiceOutput, error) {
			called = true
			if *in.ServiceName != name {
				t.Fatalf("CreateService should called with %s", name)
			}
			return nil, nil
		},
		updateService: emptyUpdateService,
	}
	return &ECS{svc: svc}, &called
}

func updateServiceShouldCalledWith(t *testing.T, name string) (*ECS, *bool) {
	called := false
	svc := &mockedECS{
		listServicesOutput: &ecssdk.ListServicesOutput{ServiceArns: testServiceArns},
		createService:      emptyCreateService,
		updateService: func(in *ecssdk.UpdateServiceInput) (*ecssdk.UpdateServiceOutput, error) {
			called = true
			if *in.Service != name {
				t.Fatalf("UpdateService should called with %s", name)
			}
			return nil, nil
		},
	}
	return &ECS{svc: svc}, &called
}

func createServiceShouldNotCalled(t *testing.T) *ECS {
	svc := &mockedECS{
		listServicesOutput: &ecssdk.ListServicesOutput{ServiceArns: testServiceArns},
		createService: func(in *ecssdk.CreateServiceInput) (*ecssdk.CreateServiceOutput, error) {
			t.Fatalf("CreateService should not called")
			return nil, nil
		},
		updateService: emptyUpdateService,
	}
	return &ECS{svc: svc}
}

func updateServiceShouldNotCalled(t *testing.T) *ECS {
	svc := &mockedECS{
		listServicesOutput: &ecssdk.ListServicesOutput{ServiceArns: testServiceArns},
		createService:      emptyCreateService,
		updateService: func(in *ecssdk.UpdateServiceInput) (*ecssdk.UpdateServiceOutput, error) {
			t.Fatalf("UpdateService should not called")
			return nil, nil
		},
	}
	return &ECS{svc: svc}
}

var emptyUpdateService = func(in *ecssdk.UpdateServiceInput) (*ecssdk.UpdateServiceOutput, error) {
	return nil, nil
}
var emptyCreateService = func(in *ecssdk.CreateServiceInput) (*ecssdk.CreateServiceOutput, error) {
	return nil, nil
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
