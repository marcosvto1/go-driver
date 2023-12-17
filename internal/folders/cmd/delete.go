package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/marcosvto1/go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func deleteFolder() *cobra.Command {
	var id int32

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a folder",
		Run: func(cmd *cobra.Command, args []string) {
			if id == 0 {
				log.Println("id are required")
				cmd.Help()
				os.Exit(1)
			}

			_, err := requests.AuthenticatedDelete(fmt.Sprintf("/folders/%d", id), nil)
			if err != nil {
				fmt.Printf("%x", err)
				os.Exit(1)
			}

			fmt.Println("Folder deleted")
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "Folder ID")

	return cmd
}
