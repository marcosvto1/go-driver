package cmd

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/marcosvto1/go-driver/internal/users"
	"github.com/marcosvto1/go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func listUser() *cobra.Command {
	var id int32

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all users",
		Run: func(cmd *cobra.Command, args []string) {
			path := "/users"
			if id > 0 {
				path = fmt.Sprintf("/users/%d", id)
			}

			res, err := requests.AuthenticatedGet(path)
			if err != nil {
				log.Fatalf("%s", err)
			}

			if id == 0 {
				var listUser []users.User
				err = json.Unmarshal(res, &listUser)
				if err != nil {
					log.Fatalf("%s", err)
				}

				t := table.NewWriter()
				t.AppendHeader(table.Row{"ID", "Name"})
				t.SetCaption("List of users")
				var rows []table.Row
				for _, u := range listUser {
					rows = append(rows, table.Row{u.ID, u.Name})
				}
				t.AppendRows(rows)

				fmt.Println(t.Render())
			} else {
				var u users.User
				err = json.Unmarshal(res, &u)
				if err != nil {
					log.Fatalf("%s", err)
				}
				// fmt.Println(Green, " - ", u.ID, " - ", u.Name)
			}
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "User ID")

	return cmd
}
