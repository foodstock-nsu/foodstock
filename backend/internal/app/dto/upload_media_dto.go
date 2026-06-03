package dto

import "io"

type UploadMediaInput struct {
	File        io.Reader
	Filename    string
	ContentType string
}

type UploadMediaOutput struct {
	MediaKey string
}
