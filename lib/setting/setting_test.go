package setting

import (
	"fmt"
	"github.com/jobtalk/pnzr/vars"
	"testing"
)

var (
	TEST_DATA_DIR = vars.TEST_DATA_DIR_ROOT + "/lib/setting"
)

func TestIsV1Setting(t *testing.T) {
	tests := []struct {
		want bool
	}{
		{false},
		{false},
		{true},
		{true},
		{false},
		{false},
	}

	for i, test := range tests {
		got := IsV1Setting(fmt.Sprintf("%s/%d.json", TEST_DATA_DIR, i))

		if got != test.want {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}
}
