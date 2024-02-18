package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/marcosvto1/go-driver/internal/folders"
	"github.com/marcosvto1/go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func getFolder() *cobra.Command {

	var id int32

	cmd := &cobra.Command{
		Use:   "get",
		Short: "get a folder",
		Run: func(cmd *cobra.Command, args []string) {
			path := "/folders/" + fmt.Sprintf("%d", id)
			if id <= 0 {
				log.Fatal("Folder id is required")
			}

			res, err := requests.AuthenticatedGet(path)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			var fc folders.FolderContent
			err = json.Unmarshal(res, &fc)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			fmt.Println("name:", fc.Folder.Name)
			fmt.Println("####### CONTENT ###########")
			if len(fc.Content) == 0 {
				fmt.Println("Folder is empty")
			}
			for _, c := range fc.Content {
				fmt.Println(c.ID, " - ", c.Type, " - ", c.Name)
			}
		},
	}

	cmd.Flags().Int32VarP(&id, "id", "", 0, "Folder id")

	return cmd
}
