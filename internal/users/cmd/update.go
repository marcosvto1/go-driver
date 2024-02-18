package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/marcosvto1/go-driver/internal/users"
	"github.com/marcosvto1/go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func updateUser() *cobra.Command {
	var id int32
	var name string

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a user",
		Run: func(cmd *cobra.Command, args []string) {
			if (id <= 0) || (name == "") {
				log.Println("id and name are required")
				fmt.Printf("\n\n### COMMAND #### \n")
				cmd.Help()
				os.Exit(1)
			}

			user := users.User{
				Name: name,
			}

			var body bytes.Buffer
			err := json.NewEncoder(&body).Encode(user)
			if err != nil {
				log.Fatalf("%v", err)
			}

			_, err = requests.AuthenticatedPut(fmt.Sprintf("/users/%d", id), &body)
			if err != nil {
				log.Fatalf("%v", err)
			}

			log.Printf("User %d updated with success", id)
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "User ID")
	cmd.Flags().StringVarP(&name, "name", "n", "", "User Name")

	return cmd
}
