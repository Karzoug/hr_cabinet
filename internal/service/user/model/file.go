package model

import "io"

type File struct {
	io.Reader
	ContentType string
	Size        int64
	Hash        string
}
