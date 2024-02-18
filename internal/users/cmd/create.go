package cmd

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	"github.com/marcosvto1/go-driver/internal/users"
	"github.com/marcosvto1/go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func createUser() *cobra.Command {
	var (
		name     string
		login    string
		password string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		Run: func(cmd *cobra.Command, args []string) {
			if (name == "") || (login == "") || (password == "") {
				log.Println("name, login and password are required")
				os.Exit(1)
			}

			u, err := users.New(name, login, password)
			if err != nil {
				log.Fatalf("%s", err)
			}

			var body bytes.Buffer
			err = json.NewEncoder(&body).Encode(u)
			if err != nil {
				log.Fatalf("%s", err)
			}

			_, err = requests.Post("/users", &body)
			if err != nil {
				log.Fatalf("%s", err)
			}

			log.Printf("User %s created with success", u.Login)
		},
	}

	cmd.Flags().StringVarP(&name, "name", "n", "", "Username")
	cmd.Flags().StringVarP(&login, "login", "l", "", "Login")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Password")

	return cmd
}
