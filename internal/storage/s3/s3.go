package s3

import (
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

const (
	bucketName = "empl-docs"
)

type File struct {
	Prefix      string
	Name        string
	ContentType string
	Size        int64
	Reader      io.Reader
}

type storage struct {
	minioClient *minio.Client
}

func New(ctx context.Context, cfg Config) (*storage, error) {
	const op = "s3 storage: new"

	endpoint := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	s := &storage{
		minioClient: minioClient,
	}

	// check to see if we already own the bucket
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		if err := minioClient.MakeBucket(ctx,
			bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		} else {
			slog.Debug(op, slog.String("bucket created", bucketName))
		}
	}

	return s, nil
}

func (s *storage) UploadFile(ctx context.Context, f File) error {
	const op = "s3 storage: upload file"

	_, err := s.minioClient.PutObject(ctx,
		bucketName,
		f.Prefix+"_"+f.Name,
		f.Reader,
		f.Size,
		minio.PutObjectOptions{ContentType: f.ContentType})
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *storage) DownloadFile(ctx context.Context, prefix, name string) (File, func() error, error) {
	const op = "s3 storage: download file"

	reader, err := s.minioClient.GetObject(ctx,
		bucketName,
		prefix+"_"+name,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return File{}, nil, fmt.Errorf("%s: %w", op, err)
	}

	info, err := reader.Stat()
	if err != nil {
		return File{}, nil, fmt.Errorf("%s: %w", op, err)
	}

	return File{
		Prefix:      prefix,
		Name:        name,
		ContentType: info.ContentType,
		Size:        info.Size,
		Reader:      reader,
	}, reader.Close, nil
}
