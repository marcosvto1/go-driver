package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/marcosvto1/go-driver/internal/folders"
	"github.com/marcosvto1/go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func createFolder() *cobra.Command {
	var name string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new folder",
		Run: func(cmd *cobra.Command, args []string) {
			if name == "" {
				log.Println("name are required")
				os.Exit(1)
				return
			}

			folder := folders.Folder{
				Name: name,
			}

			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(folder)
			if err != nil {
				log.Printf("%x", err)
				os.Exit(1)
				return
			}

			_, err = requests.AuthenticatedPost("/folders", &body)
			if err != nil {
				log.Printf("%x", err)
				os.Exit(1)
				return
			}

			log.Println("Pasta criada com sucesso")
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Folder name")

	return cmd
}
