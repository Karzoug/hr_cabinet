package s3

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/minio/minio-go/v7"

	"github.com/Employee-s-file-cabinet/backend/internal/repo/s3"
	"github.com/Employee-s-file-cabinet/backend/pkg/repoerr"
)

const (
	bucketName = "empl-docs"
)

type storage struct {
	minioClient *minio.Client
	urlExpires  time.Duration
}

func New(ctx context.Context, client *minio.Client, cfg s3.Config) (*storage, error) {
	const op = "s3 storage: new"

	s := &storage{
		minioClient: client,
		urlExpires:  cfg.URLExpires,
	}

	// check to see if we already own the bucket
	exists, err := client.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	if !exists {
		if err := client.MakeBucket(ctx,
			bucketName, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		} else {
			slog.Info(op, slog.String("bucket created", bucketName))
		}
	}

	return s, nil
}

func (s *storage) Upload(ctx context.Context, f s3.File) error {
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

func (s *storage) Download(ctx context.Context, prefix, name, etag string) (s3.File, func() error, error) {
	const op = "s3 storage: download file"

	opts := minio.GetObjectOptions{}
	if etag != "" {
		if err := opts.SetMatchETagExcept(etag); err != nil {
			return s3.File{}, nil, fmt.Errorf("%s: %w", op, err)
		}
	}
	reader, err := s.minioClient.GetObject(ctx,
		bucketName,
		prefix+"_"+name,
		opts,
	)
	if err != nil {
		return s3.File{}, nil, fmt.Errorf("%s: %w", op, err)
	}

	info, err := reader.Stat()
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		switch errResponse.StatusCode {
		case http.StatusNotFound:
			return s3.File{}, nil, fmt.Errorf("%s: %w", op, repoerr.ErrRecordNotFound)
		case http.StatusNotModified:
			return s3.File{}, nil, fmt.Errorf("%s: %w", op, repoerr.ErrRecordNotModifiedSince)
		default:
			return s3.File{}, nil, fmt.Errorf("%s: %w", op, err)
		}
	}

	return s3.File{
		Prefix:      prefix,
		Name:        name,
		ContentType: info.ContentType,
		Size:        info.Size,
		ETag:        info.ETag,
		Reader:      reader,
	}, reader.Close, nil
}

func (s *storage) PresignedURL(ctx context.Context, prefix, name string) (string, error) {
	const op = "s3 storage: object presigned url"

	url, err := s.minioClient.PresignedGetObject(ctx, bucketName, fmt.Sprintf("%s_%s", prefix, name), s.urlExpires, nil)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return url.String(), nil
}
