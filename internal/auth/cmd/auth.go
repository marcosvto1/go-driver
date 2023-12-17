package cmd

import (
	"log"

	"github.com/marcosvto1/go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func authenticate() *cobra.Command {
	var (
		user string
		pass string
	)

	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Authenticate a user",
		Run: func(cmd *cobra.Command, args []string) {
			if (user == "") || (pass == "") {
				log.Println("username and password are required")
				cmd.Help()
				return
			}

			err := requests.Auth("/auth", user, pass)
			if err != nil {
				log.Println("username and password incorrect")
				cmd.Help()
				return
			}
		},
	}

	cmd.Flags().StringVarP(&user, "user", "u", "", "Username")
	cmd.Flags().StringVarP(&pass, "pass", "p", "", "Password")

	return cmd
}
