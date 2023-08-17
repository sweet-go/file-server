package client

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"time"

	"github.com/sweet-go/file-server/model"
	"github.com/sweet-go/stdlib/helper"
	stdhttp "github.com/sweet-go/stdlib/http"
)

// ParseResponseBody parses response body from http response.
// will defer body.Close() to close the body.
// if api response is not success, it will return error.
// otherwise will parse the data to model.File.
func ParseResponseBody(body io.ReadCloser) (*model.File, error) {
	defer helper.WrapCloser(body.Close)

	var data struct {
		Response struct {
			stdhttp.StandardResponse
		} `json:"response"`
		Signature string `json:"signature"`
	}

	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}

	d, err := json.Marshal(data.Response.Data)
	if err != nil {
		return nil, err
	}

	file := &model.File{}
	if err := json.Unmarshal(d, file); err != nil {
		return nil, err
	}

	return file, nil
}

// AddDeletableMedia adds deletable media to multipart writer.
func AddDeletableMedia(writter *multipart.Writer, input UploadFileInput) error {
	if !input.IsDeletable {
		return nil
	}

	fw, err := writter.CreateFormField(model.MultipartIsDeletableKey)
	if err != nil {
		return err
	}

	_, err = fw.Write([]byte("true"))
	if err != nil {
		return err
	}

	fw, err = writter.CreateFormField(model.MultipartDeleteRuleKey)
	if err != nil {
		return err
	}

	switch input.DeleteRule {
	case model.ManualDelete:
		_, err = fw.Write([]byte(input.DeleteRule))
		if err != nil {
			return err
		}

	case model.ScheduledDelete:
		_, err := time.ParseDuration(input.DeleteAfter)
		if err != nil {
			// safety check to ensure that input.DeleteAfter is able to be parsed by time.ParseDuration
			return err
		}

		_, err = fw.Write([]byte(input.DeleteRule))
		if err != nil {
			return err
		}

		fw, err = writter.CreateFormField(model.MultipartScheduledDeleteDuration)
		if err != nil {
			return err
		}

		_, err = fw.Write([]byte(input.DeleteAfter))
		if err != nil {
			return err
		}
	}

	return nil
}
