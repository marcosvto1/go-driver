package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/marcosvto1/go-driver/internal/folders"
	"github.com/marcosvto1/go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func listFolders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List folders",
		Run: func(cmd *cobra.Command, args []string) {
			path := "/folders"

			res, err := requests.AuthenticatedGet(path)
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(1)
			}

			var fc folders.FolderContent
			err = json.Unmarshal(res, &fc)
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(1)
			}

			fmt.Println("name:", fc.Folder.Name)
			fmt.Println("######### CONTENT ##########")
			var rows []table.Row
			for _, c := range fc.Content {
				rows = append(rows, table.Row{c.ID, c.Type, c.Name})
			}

			t := table.NewWriter()
			t.AppendHeader(table.Row{"ID", "Type", "Name"})
			t.AppendRows(rows)

			fmt.Println(t.Render())

		},
	}

	return cmd
}
