package s3

import (
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type File struct {
	io.Reader
	Prefix      string
	Name        string
	ContentType string
	Size        int64
	ETag        string
}

type Client struct {
	*minio.Client
	URLExpires time.Duration
}

func NewClient(cfg Config) (Client, error) {
	endpoint := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return Client{}, err
	}

	return Client{
		Client:     mc,
		URLExpires: cfg.URLExpires,
	}, nil
}
