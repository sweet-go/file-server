// Package usecase provides usecase implementation for deletable media
package usecase

import (
	"context"
	"errors"
	"net/http"
	"path"

	"github.com/sirupsen/logrus"
	localHelper "github.com/sweet-go/file-server/internal/helper"
	"github.com/sweet-go/file-server/internal/repository"
	"github.com/sweet-go/file-server/model"
	custerr "github.com/sweet-go/stdlib/error"
	"github.com/sweet-go/stdlib/helper"
)

type deletableMedia struct {
	repo        model.DeletableMediaRepository
	storagePath string
}

// NewDeletableMediaUsecase creates new deletable media usecase
func NewDeletableMediaUsecase(repo model.DeletableMediaRepository, storagePath string) model.DeletableMediaUsecase {
	return &deletableMedia{
		repo:        repo,
		storagePath: storagePath,
	}
}

func (u *deletableMedia) DeleteMedia(ctx context.Context, input *model.DeleteMediaInput) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"input": helper.Dump(input),
		"func":  "Delete",
	})

	if err := input.Validate(); err != nil {
		return custerr.ErrChain{
			Message: "invalid input",
			Cause:   err,
			Code:    http.StatusBadRequest,
			Type:    ErrBadRequest,
		}
	}

	deletable, err := u.repo.FindByID(ctx, input.ID)
	switch err {
	default:
		logger.WithError(err).Error("failed to find deletable media")
		return &custerr.ErrChain{
			Message: "failed to find deletable media",
			Cause:   err,
			Code:    http.StatusInternalServerError,
			Type:    ErrInternal,
		}

	case repository.ErrNotFound:
		return &custerr.ErrChain{
			Message: "deletable media not found",
			Cause:   err,
			Code:    http.StatusNotFound,
			Type:    ErrNotFound,
		}

	case nil:
		break
	}

	switch deletable.DeleteRule {
	default:
		return &custerr.ErrChain{
			Message: "invalid delete rule",
			Cause:   errors.New("invalid delete rule"),
			Code:    http.StatusBadRequest,
			Type:    ErrBadRequest,
		}

	case model.ManualDelete:
		return u.handleManualDelete(ctx, deletable)
	}
}

func (u *deletableMedia) handleManualDelete(ctx context.Context, deletable *model.DeletableMedia) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"deletable": helper.Dump(deletable),
		"func":      "handleManualDelete",
	})

	filePath := path.Clean(path.Join(u.storagePath, localHelper.GenerateEncryptedFilename(deletable.Name)))
	if !localHelper.IsFileExists(filePath) {
		logger.WithError(errors.New("file not found")).Error("failed to find file")
		return &custerr.ErrChain{
			Message: "failed to find file",
			Cause:   errors.New("file not found. maybe already deleted"),
			Code:    http.StatusNotFound,
			Type:    ErrNotFound,
		}
	}

	if err := helper.DeleteFile(filePath); err != nil {
		logger.WithError(err).Error("failed to delete file")
		return &custerr.ErrChain{
			Message: "failed to delete file",
			Cause:   err,
			Code:    http.StatusInternalServerError,
			Type:    ErrInternal,
		}
	}

	if err := u.repo.Delete(ctx, deletable.ID); err != nil {
		logger.WithError(err).Error("failed to delete deletable media")
		return &custerr.ErrChain{
			Message: "failed to delete deletable media",
			Cause:   err,
			Code:    http.StatusInternalServerError,
			Type:    ErrInternal,
		}
	}

	return nil
}
