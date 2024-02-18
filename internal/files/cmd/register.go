package cmd

import "github.com/spf13/cobra"

func Register(c *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "file",
		Short: "Manage files",
	}

	cmd.AddCommand(createFile())
	cmd.AddCommand(updateFile())
	cmd.AddCommand(deleteFile())

	c.AddCommand(cmd)
}
