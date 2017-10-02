package vars

import "github.com/ieee0824/getenv"

var (
	VERSION    string
	BUILD_DATE string
	BUILD_OS   string
)
var (
	PROJECT_ROOT       = getenv.String("GOPATH") + "/src/github.com/jobtalk/pnzr"
	TEST_DATA_DIR_ROOT = PROJECT_ROOT + "/test"
)
