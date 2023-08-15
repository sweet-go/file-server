package usecase

import (
	"context"
	"errors"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/sirupsen/logrus"
	localHelper "github.com/sweet-go/file-server/internal/helper"
	"github.com/sweet-go/file-server/model"
	"github.com/sweet-go/stdlib/encryption"
	custerr "github.com/sweet-go/stdlib/error"
	"github.com/sweet-go/stdlib/helper"
)

type publicHandler struct {
	keyComponent       *encryption.KeyComponent
	storagePath        string
	deletableMediaRepo model.DeletableMediaRepository
}

// NewPublicHandler creates new public handler
func NewPublicHandler(keyComponent *encryption.KeyComponent, storagePath string, deletableMediaRepo model.DeletableMediaRepository) model.PublicHandler {
	return &publicHandler{
		keyComponent,
		storagePath,
		deletableMediaRepo,
	}
}

func (h *publicHandler) Upload(ctx context.Context, input *model.PublicUploadInput) (*model.File, error) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"input": helper.Dump(input),
		"func":  "Upload",
	})

	md, err := helper.ReadFileMetadata(input.File)
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
	if err := helper.MultipartFileSaver(input.File, path); err != nil {
		return nil, &custerr.ErrChain{
			Message: "failed to save file",
			Cause:   err,
			Code:    http.StatusInternalServerError,
			Type:    ErrInternal,
		}
	}

	// FIXME: use helper.LogIfError from stdlib to avoid linter error
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

	var deletableMedia *model.DeletableMedia
	if input.IsDeletable {
		deletableMedia = &model.DeletableMedia{
			ID:         helper.GenerateID(),
			Name:       filename,
			DeleteRule: input.DeletableMedia.DeleteRule,
		}
		if err := h.deletableMediaRepo.Create(ctx, deletableMedia); err != nil {
			// if err here, just report
			logger.WithError(err).Error("failed to write deletable media rules to db")
		}
	}

	return &model.File{
		Name:           filename,
		Size:           md.Size,
		ContentType:    md.ContentType,
		IsPublic:       true,
		CreatedAt:      time.Now(),
		DeletableMedia: deletableMedia,
	}, nil
}

func (h *publicHandler) Download(_ context.Context, filename string) (*model.File, []byte, error) {
	filepath := path.Clean(path.Join(h.storagePath, localHelper.GenerateEncryptedFilename(filename)))
	dec := localHelper.GenerateDecryptedFilename(filepath)
	opts := &encryption.FileEncryptionOpts{
		SourcePath:   filepath,
		OutputPath:   dec,
		AESKeyLength: encryption.AES128,
		Key:          h.keyComponent,
		BufferSize:   1024,
	}

	if !localHelper.IsFileExists(filepath) {
		return nil, nil, &custerr.ErrChain{
			Message: "file not found",
			Cause:   errors.New("file not found. maybe already deleted"),
			Code:    http.StatusNotFound,
			Type:    ErrNotFound,
		}
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

	// FIXME: use helper.LogIfError from stdlib to avoid linter error
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
