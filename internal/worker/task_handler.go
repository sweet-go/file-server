// Package worker contains worker server and task handler
package worker

import (
	"context"
	"errors"
	"fmt"
	"path"

	"github.com/hibiken/asynq"
	"github.com/sirupsen/logrus"
	"github.com/sweet-go/file-server/internal/config"
	localHelper "github.com/sweet-go/file-server/internal/helper"
	"github.com/sweet-go/file-server/internal/repository"
	"github.com/sweet-go/file-server/model"
	"github.com/sweet-go/stdlib/helper"
)

type th struct {
	deletableMediaRepo model.DeletableMediaRepository
}

func newTaskHandler(deletableMediaRepo model.DeletableMediaRepository) *th {
	return &th{
		deletableMediaRepo: deletableMediaRepo,
	}
}

func (t *th) handleScheduledMediaDeleteTask(ctx context.Context, task *asynq.Task) error {
	logger := logrus.WithContext(ctx).WithFields(logrus.Fields{
		"method": "HandleScheduledMediaDeleteTask",
		"task":   task,
	})

	media, err := t.deletableMediaRepo.FindForScheduledDelete(ctx, config.ScheduledDeleteBatchSize())
	switch err {
	default:
		logger.WithError(err).Error("failed to find deletable media for scheduled delete")
		return err

	case repository.ErrNotFound:
		logger.Info("no deletable media found for scheduled delete")
		return nil

	case nil:
		break
	}

	for _, m := range media {
		filePath := path.Clean(path.Join(config.StoragePath(), localHelper.GenerateEncryptedFilename(m.Name)))
		if !localHelper.IsFileExists(filePath) {
			logger.WithError(errors.New("file not found")).Error("failed to find file")
			return fmt.Errorf("file not found to be deleted: %s", m.Name)
		}

		if err := helper.DeleteFile(filePath); err != nil {
			logger.WithError(err).Error("failed to delete file")
			return err
		}

		if err := t.deletableMediaRepo.Delete(ctx, m.ID); err != nil {
			logger.WithError(err).Error("failed to delete deletable media")
			return err
		}
	}

	return nil
}
