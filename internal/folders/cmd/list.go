package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/marcosvto1/go-driver/internal/folders"
	"github.com/marcosvto1/go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func listFolders() *cobra.Command {

	var id int32

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List folders",
		Run: func(cmd *cobra.Command, args []string) {

			path := "/folders"
			if id > 0 {
				path = fmt.Sprintf("/folders/%d", id)
			}

			res, err := requests.AuthenticatedGet(path)
			if err != nil {
				fmt.Printf("%x", err)
				os.Exit(1)
			}

			var fc folders.FolderContent
			err = json.Unmarshal(res, &fc)
			if err != nil {
				fmt.Printf("%x", err)
				os.Exit(1)
			}

			fmt.Println(fc.Folder.Name)
			fmt.Println("##################")
			for _, c := range fc.Content {
				fmt.Println(c.ID, " - ", c.Type, " - ", c.Name)
			}
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "Folder id")

	return cmd
}
