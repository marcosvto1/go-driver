package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/marcosvto1/go-driver/internal/files"
	"github.com/marcosvto1/go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func updateFile() *cobra.Command {
	var id int32
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a file",
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

			path := fmt.Sprintf("/files/%d", id)

			file := files.File{
				Name: name,
			}

			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(file)
			if err != nil {
				fmt.Printf("%x", err)
				cmd.Help()
				os.Exit(1)
			}

			_, err = requests.AuthenticatedPut(path, &body)
			if err != nil {
				fmt.Printf("%x", err)
				cmd.Help()
				os.Exit(1)
			}

			fmt.Println("file updated")
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "file id")
	cmd.Flags().StringVarP(&name, "name", "n", "", "File name")

	return cmd
}
