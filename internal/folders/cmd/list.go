package cmd

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/marcosvto1/go-driver/internal/folders"
	"github.com/marcosvto1/go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func listFolders() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List folders",
		Run: func(cmd *cobra.Command, args []string) {
			res, err := requests.AuthenticatedGet("/folders")
			if err != nil {
				fmt.Printf("%x", err)
			}

			data, err := io.ReadAll(res.Body)
			if err != nil {
				fmt.Printf("%x", err)
			}

			var listFolders []folders.Folder
			json.Unmarshal(data, &listFolders)

			fmt.Println(listFolders)
		},
	}

	return cmd
}
