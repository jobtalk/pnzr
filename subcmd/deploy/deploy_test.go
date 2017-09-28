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
