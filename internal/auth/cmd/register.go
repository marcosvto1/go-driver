package cmd

import "github.com/spf13/cobra"

func Register(c *cobra.Command) {
	// TODO
	c.AddCommand(authenticate())
}
