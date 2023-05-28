package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/marcosvto1/go-driver/internal/bucket"
	"github.com/marcosvto1/go-driver/internal/queue"
)

func main() {
	// RabbitMQ config
	qcfg := queue.RabbitMQConfig{
		URL:       os.Getenv("RABBIT_URL"),
		TopicName: os.Getenv("RABBIT_TOPIC_NAME"),
		Timeout:   time.Now().Add(30 * time.Second),
	}

	qc, err := queue.New(queue.RabbitMQ, qcfg)
	if err != nil {
		panic(err)
	}

	c := make(chan queue.QueueDTO)

	qc.Consume(c)

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
		panic(err)
	}

	for msg := range c {
		src := fmt.Sprintf("%s/%s", msg.Path, msg.Filename)
		dst := fmt.Sprintf("%d_%s", msg.ID, msg.Filename)
		file, err := b.Download(src, dst)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		body, err := io.ReadAll(file)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		var buf bytes.Buffer
		zw := gzip.NewWriter(&buf)
		_, err = zw.Write(body)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		if err := zw.Close(); err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		zr, err := gzip.NewReader(&buf)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		err = b.Upload(zr, src)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		err = os.Remove(dst)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}
	}
}
