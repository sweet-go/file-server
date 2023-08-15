package model

import (
	"context"
	"mime/multipart"
	"time"
)

const (
	MultipartFileKey        = "file"
	MultipartIsDeletableKey = "is_deletable"
	MultipartDeleteRuleKey  = "delete_rule"
)

type File struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`

	DeletableMedia *DeletableMedia `json:"deletable_media,omitempty"`
}

type PublicUploadInput struct {
	File           *multipart.FileHeader
	IsDeletable    bool
	DeletableMedia *DeletableMedia
}

type PublicHandler interface {
	Upload(ctx context.Context, input *PublicUploadInput) (*File, error)
	Download(ctx context.Context, filename string) (*File, []byte, error)
}
