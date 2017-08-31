package vault

type VaultCommand struct {
}

func (c *VaultCommand) Help() string {
	c.Help()
}

func (v *VaultCommand) Run(args []string) int {

	return 0
}

func (c *VaultCommand) Synopsis() string {
	return c.Help()
}
