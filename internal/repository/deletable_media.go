// Package repository provides the repository layer for the file server service
package repository

import (
	"context"

	"github.com/sirupsen/logrus"
	"github.com/sweet-go/file-server/model"
	"gorm.io/gorm"
)

type deletableMediaRepo struct {
	db *gorm.DB
}

// NewDeletableMediaRepository creates a new deletable media repository
func NewDeletableMediaRepository(db *gorm.DB) model.DeletableMediaRepository {
	return &deletableMediaRepo{
		db: db,
	}
}

func (r *deletableMediaRepo) Create(ctx context.Context, input *model.DeletableMedia) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"method": "Create",
		"input":  input,
	})

	if err := r.db.WithContext(ctx).Create(input).Error; err != nil {
		logger.WithError(err).Error("failed to create deletable media")
		return err
	}

	return nil
}

func (r *deletableMediaRepo) FindByID(ctx context.Context, id string) (*model.DeletableMedia, error) {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"method": "FindByID",
		"id":     id,
	})

	var deletableMedia *model.DeletableMedia
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&deletableMedia).Error
	switch err {
	default:
		logger.WithError(err).Error("failed to find deletable media")
		return nil, err

	case gorm.ErrRecordNotFound:
		return nil, ErrNotFound

	case nil:
		return deletableMedia, nil
	}
}

func (r *deletableMediaRepo) Delete(ctx context.Context, id string) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"method": "Delete",
		"id":     id,
	})

	if err := r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.DeletableMedia{}).Error; err != nil {
		logger.WithError(err).Error("failed to delete deletable media")
		return err
	}

	return nil
}
