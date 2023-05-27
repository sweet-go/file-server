package usecase

import (
	"context"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"time"

	localHelper "github.com/sweet-go/file-server/internal/helper"
	"github.com/sweet-go/file-server/model"
	"github.com/sweet-go/stdlib/encryption"
	custerr "github.com/sweet-go/stdlib/error"
	"github.com/sweet-go/stdlib/helper"
)

type publicHandler struct {
	keyComponent *encryption.KeyComponent
	storagePath  string
}

func NewPublicHandler(keyComponent *encryption.KeyComponent, storagePath string) model.PublicHandler {
	return &publicHandler{
		keyComponent,
		storagePath,
	}
}

func (h *publicHandler) Upload(ctx context.Context, file *multipart.FileHeader) (*model.File, error) {
	md, err := helper.ReadFileMetadata(file)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "failed to read file metadata",
			Cause:   err,
			Code:    http.StatusBadRequest,
			Type:    ErrBadRequest,
		}
	}

	filename := localHelper.GenerateUploadedFilename(md.Name)
	path := path.Clean(path.Join(h.storagePath, filename))
	if err := helper.MultipartFileSaver(file, path); err != nil {
		return nil, &custerr.ErrChain{
			Message: "failed to save file",
			Cause:   err,
			Code:    http.StatusInternalServerError,
			Type:    ErrInternal,
		}
	}

	defer helper.DeleteFile(path)

	encOpts := &encryption.FileEncryptionOpts{
		SourcePath:   path,
		OutputPath:   localHelper.GenerateEncryptedFilename(path),
		AESKeyLength: encryption.AES128,
		Key:          h.keyComponent,
		BufferSize:   1024,
	}

	_, err = encryption.EncryptFile(encOpts)
	if err != nil {
		return nil, &custerr.ErrChain{
			Message: "failed to encrypt file",
			Cause:   err,
			Code:    http.StatusInternalServerError,
			Type:    ErrInternal,
		}
	}

	return &model.File{
		Name:        filename,
		Size:        md.Size,
		ContentType: md.ContentType,
		IsPublic:    true,
		CreatedAt:   time.Now(),
	}, nil
}

func (h *publicHandler) Download(ctx context.Context, filename string) (*model.File, []byte, error) {
	filepath := path.Clean(path.Join(h.storagePath, localHelper.GenerateEncryptedFilename(filename)))
	dec := localHelper.GenerateDecryptedFilename(filepath)
	opts := &encryption.FileEncryptionOpts{
		SourcePath:   filepath,
		OutputPath:   dec,
		AESKeyLength: encryption.AES128,
		Key:          h.keyComponent,
		BufferSize:   1024,
	}

	err := encryption.DecryptFile(opts)
	if err != nil {
		return nil, nil, &custerr.ErrChain{
			Message: "failed to decrypt file",
			Cause:   err,
			Code:    http.StatusInternalServerError,
			Type:    ErrInternal,
		}
	}

	defer helper.DeleteFile(dec)

	f, err := os.ReadFile(dec)
	if err != nil {
		return nil, nil, &custerr.ErrChain{
			Message: "failed to read file",
			Cause:   err,
			Code:    http.StatusInternalServerError,
			Type:    ErrInternal,
		}
	}

	s, err := os.Stat(dec)
	if err != nil {
		return nil, nil, &custerr.ErrChain{
			Message: "failed to read file metadata",
			Cause:   err,
			Code:    http.StatusInternalServerError,
			Type:    ErrInternal,
		}
	}

	return &model.File{
		Name:        filename,
		Size:        s.Size(),
		ContentType: http.DetectContentType(f),
	}, f, nil
}
