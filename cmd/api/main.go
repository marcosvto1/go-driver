package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/go-chi/chi/v5"
	"github.com/marcosvto1/go-driver/internal/auth"
	"github.com/marcosvto1/go-driver/internal/bucket"
	"github.com/marcosvto1/go-driver/internal/files"
	"github.com/marcosvto1/go-driver/internal/folders"
	"github.com/marcosvto1/go-driver/internal/queue"
	"github.com/marcosvto1/go-driver/internal/users"
	"github.com/marcosvto1/go-driver/pkg/database"
)

func main() {
	db, b, qc := getSessions()
	defer db.Close()

	// Create new Router using framework CHI
	r := chi.NewRouter()

	// Define endpoint
	r.Post("/auth", auth.HandleAuth(func(login, pass string) (auth.Authenticated, error) {
		return users.Authenticate(login, pass)
	}))
	files.SetRoutes(r, db, b, qc)
	folders.SetRoute(r, db)
	users.SetRoutes(r, db)

	err := http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APPLICATION_PORT")), r)
	if err != nil {
		log.Fatal(err)
	}
}

func getSessions() (*sql.DB, *bucket.Bucket, *queue.Queue) {
	// Database Config
	db, err := database.NewConnection()
	if err != nil {
		log.Fatal(err)
	}

	// RabbitMQ config
	qcfg := queue.RabbitMQConfig{
		URL:       os.Getenv("RABBIT_URL"),
		TopicName: os.Getenv("RABBIT_TOPIC_NAME"),
		Timeout:   time.Now().Add(30 * time.Second),
	}

	qc, err := queue.New(queue.RabbitMQ, qcfg)
	if err != nil {
		log.Fatal(err)
	}

	// Bucket AWS Provider Config
	bcfg := bucket.AwsProviderConfig{
		Config: &aws.Config{
			Region:      aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_KEY"), os.Getenv("AWS_SECRET"), ""),
		},
		BucketDownload: "drive-raw",
		BucketUpload:   "drive-compact-gzip",
	}

	b, err := bucket.New(bucket.AwsProvider, bcfg)
	if err != nil {
		log.Fatal(err)
	}

	return db, b, qc
}
