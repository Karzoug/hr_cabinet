package s3

import "io"

type UploadableFile struct {
	Prefix      string
	Name        string
	ContentType string
	Size        int64
	Reader      io.Reader
}
