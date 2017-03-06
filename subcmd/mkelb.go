package subcmd

import (
	"encoding/json"

	"github.com/ieee0824/thor/mkelb"
)

type mkelbConfigure struct {
	*mkelb.Setting
}

type mkelbParam struct {
}

func (p mkelbParam) String() string {
	bin, err := json.Marshal(p)
	if err != nil {
		return ""
	}
	return string(bin)
}

type MkELB struct{}

func (c *MkELB) Help() string {
	return ""
}

func (c *MkELB) Run(args []string) int {
	return 0
}

func (c *MkELB) Synopsis() string {
	return ""
}
