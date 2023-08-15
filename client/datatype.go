package client

import (
	"context"

	"github.com/sweet-go/file-server/model"
)

// UploadFileInput is the input for Upload method
type UploadFileInput struct {
	FullPath       string
	IsDeletable    bool
	DeletableMedia *model.DeletableMedia
}

type DeleteFileInput struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Utility is the interface for client utility
type Utility interface {
	Upload(ctx context.Context, input UploadFileInput) (*model.File, error)
	Delete(ctx context.Context, input *DeleteFileInput) error
}
