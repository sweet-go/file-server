package model

import (
	"context"
	"mime/multipart"
	"time"
)

type File struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
}

type PublicHandler interface {
	Upload(ctx context.Context, file *multipart.FileHeader) (*File, error)
	Download(ctx context.Context, filename string) (*File, []byte, error)
}
