package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	authCMD "github.com/marcosvto1/go-driver/internal/auth/cmd"
	fileCMD "github.com/marcosvto1/go-driver/internal/files/cmd"
	folderCMD "github.com/marcosvto1/go-driver/internal/folders/cmd"
	userCMD "github.com/marcosvto1/go-driver/internal/users/cmd"
)

var RootCMD = &cobra.Command{}

func main() {
	godotenv.Load("../../.env")

	authCMD.Register(RootCMD)
	userCMD.Register(RootCMD)
	fileCMD.Register(RootCMD)
	folderCMD.Register(RootCMD)

	if err := RootCMD.Execute(); err != nil {
		log.Fatal(err)
	}
}
