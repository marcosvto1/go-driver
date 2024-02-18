package cmd

import "github.com/spf13/cobra"

func Register(c *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Manage users",
	}

	cmd.AddCommand(createUser())
	cmd.AddCommand(updateUser())
	cmd.AddCommand(listUser())
	cmd.AddCommand(deleteUser())

	c.AddCommand(cmd)
}
