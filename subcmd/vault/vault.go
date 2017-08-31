package vault

import (
	"fmt"
	"github.com/jobtalk/pnzr/subcmd/vault/decrypt"
	"github.com/jobtalk/pnzr/subcmd/vault/edit"
	"github.com/jobtalk/pnzr/subcmd/vault/encrypt"
	"github.com/jobtalk/pnzr/subcmd/vault/view"
	"github.com/jobtalk/pnzr/vars"
	"github.com/mitchellh/cli"
)

type VaultCommand struct {
	cli  *cli.CLI
	args []string
}

func New(args []string) *VaultCommand {
	ret := &VaultCommand{
		cli:  cli.NewCLI("vault", vars.VERSION),
		args: args,
	}

	ret.cli.Commands = map[string]cli.CommandFactory{
		"edit": func() (cli.Command, error) {
			return &edit.EditCommand{}, nil
		},
		"view": func() (cli.Command, error) {
			return &view.ViewCommand{}, nil
		},
		"encrypt": func() (cli.Command, error) {
			return &encrypt.EncryptCommand{}, nil
		},
		"decrypt": func() (cli.Command, error) {
			return &decrypt.DecryptCommand{}, nil
		},
	}

	return ret
}

func (v *VaultCommand) Run(args []string) int {
	v.cli.Args = args

	exitCode, err := v.cli.Run()
	if err != nil {
		panic(err)
	}

	return exitCode
}

func (v *VaultCommand) Help() string {
	if 2 <= len(v.args) {
		v.Run(v.args[1:])
		return ""
	}

	return v.Synopsis()
}

func (v *VaultCommand) Synopsis() string {
	var help = " (subcmd)\n"
	help += "    subcmd list:\n"

	for v, _ := range v.cli.Commands {
		help += fmt.Sprintf("        %s\n", v)
	}

	return help
}
