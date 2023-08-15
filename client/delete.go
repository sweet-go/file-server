package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/sirupsen/logrus"
	"github.com/sweet-go/file-server/model"
	"github.com/sweet-go/stdlib/helper"
)

func (c *impl) Delete(ctx context.Context, input *DeleteFileInput) error {
	logger := logrus.WithContext(ctx).WithFields((logrus.Fields{
		"method": "Delete",
		"input":  helper.Dump(input),
	}))

	body := &model.DeleteMediaInput{
		ID: input.ID,
	}

	if err := body.Validate(); err != nil {
		logger.WithError(err).Error("failed to validate input")
		return err
	}

	b, err := json.Marshal(body)
	if err != nil {
		logger.WithError(err).Error("failed to marshal input")
		return err
	}

	url := c.baseURL + "/public/"
	req, err := http.NewRequest(http.MethodDelete, url, bytes.NewReader(b))
	if err != nil {
		logger.WithError(err).Error("failed to create request")
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpclient.Do(req)
	if err != nil {
		logger.WithError(err).Error("failed to do request")
		return err
	}

	defer helper.WrapCloser(resp.Body.Close)

	if resp.StatusCode != http.StatusOK {
		logger.WithField("status_code", resp.StatusCode).Error("failed to delete file")
		return errors.New("file server return non 200 status when delete file")
	}

	return nil
}
