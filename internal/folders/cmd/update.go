package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/marcosvto1/go-driver/internal/folders"
	"github.com/marcosvto1/go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func updateFolder() *cobra.Command {
	var (
		id   int32
		name string
	)

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a folder",
		Run: func(cmd *cobra.Command, args []string) {
			if id == 0 {
				log.Println("id are required")
				cmd.Help()
				os.Exit(1)
			}

			if name == "" {
				log.Println("name are required")
				cmd.Help()
				os.Exit(1)
			}

			path := fmt.Sprintf("/folders/%d", id)

			folder := folders.Folder{
				Name: name,
			}

			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(folder)
			if err != nil {
				log.Printf("%x", err)

				cmd.Help()
				os.Exit(1)
			}

			_, err = requests.AuthenticatedPut(path, &body)
			if err != nil {
				log.Printf("%x", err)

				cmd.Help()
				os.Exit(1)
			}
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "Folder id")
	cmd.Flags().StringVarP(&name, "name", "n", "", "Folder Name")

	return cmd
}
