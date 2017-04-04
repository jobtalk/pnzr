package main

import (
	"log"
	"os"

	"github.com/jobtalk/thor/subcmd/mkelb"
	"github.com/jobtalk/thor/subcmd/vault"
	"github.com/mitchellh/cli"
)

const (
	VERSION = "0.01"
)

func init() {
	log.SetFlags(log.Llongfile)
}

func main() {
	c := cli.NewCLI("thor", VERSION)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		/*
			"deploy": func() (cli.Command, error) {
				return &deploy.Deploy{}, nil
			},
		*/
		"mkelb": func() (cli.Command, error) {
			return &mkelb.MkELB{}, nil
		},
		"vault": func() (cli.Command, error) {
			return &vault.Vault{}, nil
		},
	}
	exitCode, err := c.Run()
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(exitCode)
}
