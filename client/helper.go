package client

import (
	"encoding/json"
	"io"

	"github.com/sweet-go/file-server/model"
	"github.com/sweet-go/stdlib/helper"
)

// ParseResponseBody parses response body from http response.
// will defer body.Close() to close the body.
// if api response is not success, it will return error.
// otherwise will parse the data to model.File.
func ParseResponseBody(body io.ReadCloser) (*model.File, error) {
	defer helper.WrapCloser(body.Close)

	var data struct {
		Response struct {
			Data *model.File
		} `json:"response"`
		Signature string
	}
	if err := json.NewDecoder(body).Decode(&data); err != nil {
		return nil, err
	}

	return data.Response.Data, nil
}
