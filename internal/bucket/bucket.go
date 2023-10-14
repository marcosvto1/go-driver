package bucket

import (
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
)

const (
	AwsProvider BucketType = iota
	MockProvider
)

type BucketType int

type BucketInterface interface {
	Upload(reader io.Reader, destiny string) error
	Download(source string, destiny string) (*os.File, error)
	Delete(string) error
}

type Bucket struct {
	provider BucketInterface
}

func New(bt BucketType, cfg any) (b *Bucket, err error) {
	b = new(Bucket)

	providerConfig := reflect.TypeOf(cfg)
	switch bt {
	case AwsProvider:
		if providerConfig.Name() != "AwsProviderConfig" {
			return nil, fmt.Errorf("config need's to be of type AwsProviderConfig")
		}

		b.provider = newAwsProvider(cfg.(AwsProviderConfig))
	case MockProvider:
		cf := cfg.(MockBucketConfig)
		b.provider = &MockBucket{
			content:     make(map[string][]byte),
			mockOptions: cf,
		}
	default:
		log.Fatal("type not implemented")
	}
	return
}

func (b *Bucket) Upload(reader io.Reader, destiny string) error {
	return b.provider.Upload(reader, destiny)
}

func (b *Bucket) Download(source, destiny string) (*os.File, error) {
	return b.provider.Download(source, destiny)
}

func (b *Bucket) Delete(destiny string) error {
	return b.provider.Delete(destiny)
}
