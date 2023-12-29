package s3

import (
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type File struct {
	Prefix      string
	Name        string
	ContentType string
	Size        int64
	Reader      io.Reader
}

func NewClient(cfg Config) (*minio.Client, error) {
	endpoint := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
}
