package main

import (
	"log"
	"os"

	"github.com/ieee0824/thor/subcmd"
	"github.com/mitchellh/cli"
)

const (
	VERSION = "0.01"
)

func main() {
	log.SetFlags(log.Llongfile)
	/*
		awsConfig := &aws.Config{
			Region: aws.String("ap-northeast-1"),
		}
	*/

	c := cli.NewCLI("thor", VERSION)
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"deploy": func() (cli.Command, error) {
			return &subcmd.Deploy{}, nil
		},
	}
	exitCode, err := c.Run()
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(exitCode)
}
