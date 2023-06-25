package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/sweet-go/file-server/model"
	"github.com/sweet-go/stdlib/helper"
)

func (c *impl) Upload(ctx context.Context, input UploadFileInput) (*model.File, error) {
	logger := logrus.WithFields(logrus.Fields{
		"method": "Upload",
		"input":  helper.Dump(input),
	})

	cleanPath := path.Clean(input.FullPath)
	file, err := os.Open(cleanPath)
	if err != nil {
		logger.WithError(err).Error("failed to open file")
		return nil, err
	}

	defer helper.WrapCloser(file.Close)

	body := &bytes.Buffer{}
	writter := multipart.NewWriter(body)
	part, err := writter.CreateFormFile(model.MultipartFileKey, file.Name())
	if err != nil {
		logger.WithError(err).Error("failed to create form file")
		return nil, err
	}

	_, err = io.Copy(part, file)
	if err != nil {
		logger.WithError(err).Error("failed to copy file")
		return nil, err
	}

	helper.WrapCloser(writter.Close)

	url := fmt.Sprintf("%s/upload/", c.baseURL)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		logger.WithError(err).Error("failed to create request")
		return nil, err
	}

	req.Header.Add("Content-Type", writter.FormDataContentType())
	resp, err := c.httpclient.Do(req)
	if err != nil {
		logger.WithError(err).Error("failed to do request")
		return nil, err
	}

	if resp.StatusCode == http.StatusInternalServerError {
		logger.WithError(err).Error("failed to upload file because file server return internal server error")
		return nil, err
	}

	return ParseResponseBody(resp.Body)
}
