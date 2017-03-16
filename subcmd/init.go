package subcmd

import (
	"log"

	i "github.com/ieee0824/thor/init"
	termbox "github.com/nsf/termbox-go"
)

type Init struct{}

func (c *Init) Run(args []string) int {
	err := termbox.Init()
	if err != nil {
		log.Fatalln(err)
	}
	defer termbox.Close()
	i.RunInit()
	return 0
}

func (c *Init) Help() string {
	return ""
}

func (c *Init) Synopsis() string {
	return ""
}
