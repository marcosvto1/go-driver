package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"

	"github.com/marcosvto1/go-driver/pkg/requests"
	"github.com/spf13/cobra"
)

func createFile() *cobra.Command {
	var (
		folderId int32
		filename string
	)

	cmd := &cobra.Command{
		Use:   "upload",
		Short: "upload a new file",
		Run: func(cmd *cobra.Command, args []string) {
			if filename == "" {
				log.Println("filename are required")
				os.Exit(1)
			}

			file, err := os.Open(filename)
			if err != nil {
				log.Printf("%v", err)
				os.Exit(1)
			}

			defer file.Close()

			var body bytes.Buffer

			mw := multipart.NewWriter(&body)
			w, err := mw.CreateFormFile("file", filepath.Base(file.Name()))
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(1)
			}
			io.Copy(w, file)

			if folderId > 0 {
				w, err := mw.CreateFormField("folder_id")
				if err != nil {
					fmt.Printf("%v", err)
					os.Exit(1)
				}

				_, err = w.Write([]byte(strconv.Itoa(int(folderId))))
				if err != nil {
					fmt.Printf("%v", err)
					os.Exit(1)
				}
			}

			mw.Close()

			headers := map[string]string{
				"Content-Type": mw.FormDataContentType(),
			}

			_, err = requests.AuthenticatedWithHeaders("/files", &body, headers)
			if err != nil {
				fmt.Printf("%v", err)
				os.Exit(1)
			}

			log.Println("File created")
		},
	}

	cmd.Flags().Int32VarP(&folderId, "folder-id", "i", 0, "Folder ID")
	cmd.Flags().StringVarP(&filename, "filename", "f", "", "File name")

	return cmd
}
