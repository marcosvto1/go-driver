package cmd

import "github.com/spf13/cobra"

func Register(c *cobra.Command) {
	cmd := &cobra.Command{
		Use:   "folder",
		Short: "Manage folders",
	}

	cmd.AddCommand(createFolder())
	cmd.AddCommand(listFolders())
	cmd.AddCommand(updateFolder())
	cmd.AddCommand(deleteFolder())
	cmd.AddCommand(getFolder())

	c.AddCommand(cmd)
}
