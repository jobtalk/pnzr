package deploy

import "os"

var (
	TEST_DIR = os.Getenv("GOPATH") + "/src/github.com/jobtalk/pnzr/test"
)

func init() {
}
