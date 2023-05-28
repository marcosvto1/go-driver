package bucket

import (
	"io"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
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
	return
}

func (p *awsProvider) Upload(reader io.Reader, destinyKey string) error {
	return nil
}
func (p *awsProvider) Delete(destinyKey string) error {
	return nil
}
