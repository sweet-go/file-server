package model

import (
	"context"
	"mime/multipart"
	"time"
)

// list constant for easier usage between internal func and client package to lookup certain input
const (
	MultipartFileKey        = "file"
	MultipartIsDeletableKey = "is_deletable"
	MultipartDeleteRuleKey  = "delete_rule"
)

// File is a model for file
type File struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`

	DeletableMedia *DeletableMedia `json:"deletable_media,omitempty"`
}

// PublicUploadInput is input to upload file
type PublicUploadInput struct {
	File           *multipart.FileHeader
	IsDeletable    bool
	DeletableMedia *DeletableMedia
}

// PublicHandler is an interface for public handler
type PublicHandler interface {
	Upload(ctx context.Context, input *PublicUploadInput) (*File, error)
	Download(ctx context.Context, filename string) (*File, []byte, error)
}
