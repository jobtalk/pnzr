package deploy

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ecs"
)

var (
	TEST_DIR = os.Getenv("GOPATH") + "/src/github.com/jobtalk/pnzr/test"
)

func init() {
}

func TestProgressNewRun(t *testing.T) {

	deployments := []*ecs.Deployment{
		{DesiredCount: aws.Int64(3), RunningCount: aws.Int64(3)},
		{DesiredCount: aws.Int64(0), RunningCount: aws.Int64(3)},
		{DesiredCount: aws.Int64(3), RunningCount: aws.Int64(0)},
	}
	tests := []struct {
		revision  int
		index     int
		pRevision int
		want      bool
	}{
		{
			10,
			0,
			10,
			true,
		},
		{
			10,
			1,
			10,
			false,
		},
		{
			10,
			2,
			10,
			false,
		},
		{
			11,
			0,
			10,
			false,
		},
		{
			11,
			1,
			10,
			false,
		},
		{
			11,
			2,
			10,
			false,
		},
		{
			10,
			0,
			11,
			false,
		},
		{
			10,
			1,
			11,
			false,
		},
		{
			10,
			2,
			11,
			false,
		},
	}

	for _, test := range tests {
		p := &Progress{revision: test.pRevision}
		got := p.progressNewRun(test.revision, test.index, deployments)
		if got != test.want {
			t.Fatalf("want %v, but %v:", test.want, got)
		}
	}
	deployments = deployments[:0]
	for _, test := range tests {
		test.want = false
		p := &Progress{revision: test.pRevision}
		got := p.progressNewRun(test.revision, test.index, deployments)
		if got != test.want {
			t.Fatalf("want %v, but %v:", test.want, got)
		}
	}
}

func TestProgressOldStop(t *testing.T) {
	deployments := []*ecs.Deployment{
		{DesiredCount: aws.Int64(3)},
	}
	tests := []struct {
		revision  int
		pRevision int
		want      bool
	}{
		{
			10,
			10,
			true,
		},
		{
			11,
			10,
			false,
		},
		{
			10,
			11,
			false,
		},
	}
	for _, test := range tests {
		p := &Progress{
			revision: test.pRevision,
		}
		got := p.progressOldStop(test.revision, 0, deployments)
		if got != test.want {
			t.Fatalf("want %v, but %v:", test.want, got)
		}
	}
	deployments = append(deployments, &ecs.Deployment{DesiredCount: aws.Int64(0)})
	for _, test := range tests {
		test.want = false
		p := &Progress{
			revision: test.pRevision,
		}
		got := p.progressOldStop(test.revision, 0, deployments)
		if got != test.want {
			t.Fatalf("want %v, but %v:", test.want, got)
		}
	}

}

func TestGetNextState(t *testing.T) {
	deployments1 := Deployments{
		{
			TaskDefinition: aws.String(":10"),
			DesiredCount:   aws.Int64(3),
			RunningCount:   aws.Int64(3),
		},
		{
			TaskDefinition: aws.String(":11"),
			DesiredCount:   aws.Int64(3),
			RunningCount:   aws.Int64(3),
		},
	}
	deployments2 := []*ecs.Deployment{
		{
			TaskDefinition: aws.String(":11"),
			DesiredCount:   aws.Int64(3),
			RunningCount:   aws.Int64(3),
		},
	}
	tests := []struct {
		state       string
		pRevision   int
		deployments Deployments
		wantState   string
		wantMessage string
	}{
		{
			"initial",
			11,
			deployments1,
			"launched",
			"(2/3) デプロイ対象のコンテナが全て起動しました",
		},
		{
			"initial",
			10,
			deployments1,
			"error",
			"正常な処理が行われませんでした。",
		},
		{
			"launched",
			11,
			deployments2,
			"done",
			"(3/3) 古いコンテナが全て停止しました",
		},
		{
			"launched",
			11,
			deployments1,
			"launched",
			"(2/3) デプロイ対象のコンテナが全て起動しました",
		},
		{
			"launched",
			11,
			deployments1,
			"launched",
			"(2/3) デプロイ対象のコンテナが全て起動しました",
		},
		{
			"launched",
			10,
			deployments1,
			"error",
			"正常な処理が行われませんでした。",
		},
		{
			"",
			11,
			deployments1,
			"error",
			"正常な処理が行われませんでした。",
		},
		{
			"",
			11,
			deployments2,
			"error",
			"正常な処理が行われませんでした。",
		},
	}

	for i, test := range tests {
		p := &Progress{
			revision: test.pRevision,
		}
		gotState, gotMessage := p.getNextState(test.state, test.deployments)
		if gotState != test.wantState || gotMessage != test.wantMessage {
			t.Fatalf("index %d :want %s, %s, but %s, %s:", i, test.wantState, test.wantMessage, gotState, gotMessage)
		}
	}

}
