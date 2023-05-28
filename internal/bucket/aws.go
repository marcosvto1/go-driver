package bucket

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

type AwsProviderConfig struct {
	Config         *aws.Config
	BucketDownload string
	BucketUpload   string
}

type awsProvider struct {
	session        *session.Session
	bucketDownload string
	bucketUpload   string
}

func newAwsProvider(cfg AwsProviderConfig) *awsProvider {
	c, _ := session.NewSession(cfg.Config)

	return &awsProvider{
		session:        c,
		bucketDownload: cfg.BucketDownload,
		bucketUpload:   cfg.BucketUpload,
	}
}

func (p *awsProvider) Download(source, destiny string) (file *os.File, err error) {
	file, err = os.Create(destiny)
	if err != nil {
		return
	}
	defer file.Close()

	downloader := s3manager.NewDownloader(p.session)

	_, err = downloader.Download(file, &s3.GetObjectInput{
		Bucket: aws.String(p.bucketDownload),
		Key:    aws.String(source),
	})

	return
}

func (p *awsProvider) Upload(reader io.Reader, destinyKey string) error {
	uploader := s3manager.NewUploader(p.session)

	_, err := uploader.Upload(&s3manager.UploadInput{
		Key:    aws.String(destinyKey),
		Bucket: aws.String(p.bucketUpload),
		Body:   reader,
	})

	if err != nil {
		return err
	}

	return nil
}

func (p *awsProvider) Delete(destinyKey string) error {
	svc := s3.New(p.session)

	_, err := svc.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(p.bucketDownload),
		Key:    aws.String(destinyKey),
	})

	if err != nil {
		return err
	}

	return svc.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(p.bucketDownload),
		Key:    aws.String(destinyKey),
	})
}
