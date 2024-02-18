package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/marcosvto1/go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func deleteUser() *cobra.Command {
	var id int32

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a user",
		Run: func(cmd *cobra.Command, args []string) {
			if id == 0 {
				log.Println("id are required")
				os.Exit(1)
			}

			_, err := requests.AuthenticatedDelete(fmt.Sprintf("/users/%d", id), nil)
			if err != nil {
				log.Printf("%s", err)
				os.Exit(1)
			}

			fmt.Println("User deleted")
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "User ID")

	return cmd
}
