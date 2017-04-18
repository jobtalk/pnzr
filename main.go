package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/jobtalk/thor/subcmd/deploy"
	"github.com/jobtalk/thor/subcmd/mkelb"
	"github.com/jobtalk/thor/subcmd/update"
	"github.com/jobtalk/thor/subcmd/vault"
	"github.com/jobtalk/thor/vars"
	"github.com/mitchellh/cli"
)

var (
	VERSION    string
	BUILD_DATE string
	BUILD_OS   string
)

func generateBuildInfo() string {
	ret := fmt.Sprintf("Build version: %s\n", VERSION)
	ret += fmt.Sprintf("Go version: %s\n", runtime.Version())
	ret += fmt.Sprintf("Build Date: %s\n", BUILD_DATE)
	ret += fmt.Sprintf("Build OS: %s\n", BUILD_OS)
	return ret
}

func init() {
	if VERSION == "" {
		VERSION = "unknown"
	}
	VERSION = generateBuildInfo()
	log.SetFlags(log.Llongfile)

	vars.VERSION = &VERSION
	vars.BUILD_DATE = &BUILD_DATE
	vars.BUILD_OS = &BUILD_OS
}

func main() {
	c := cli.NewCLI("thor", VERSION)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"deploy": func() (cli.Command, error) {
			return &deploy.Deploy{}, nil
		},
		"mkelb": func() (cli.Command, error) {
			return &mkelb.MkELB{}, nil
		},
		"vault": func() (cli.Command, error) {
			return &vault.Vault{}, nil
		},
		"update": func() (cli.Command, error) {
			return &update.Update{}, nil
		},
	}
	exitCode, err := c.Run()
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(exitCode)
}
