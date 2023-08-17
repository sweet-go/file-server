package httpsvc

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sweet-go/file-server/model"
	stdhttp "github.com/sweet-go/stdlib/http"
)

func (s *Service) handlePublicUpload() echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile(model.MultipartFileKey)
		if err != nil {
			return s.apiRespGenerator.GenerateEchoAPIResponse(c, &stdhttp.StandardResponse{
				Success: false,
				Message: "missing required file",
				Status:  http.StatusBadRequest,
				Error:   echo.ErrBadRequest.Error(),
			}, nil)
		}

		isDeleteable := false
		var deletable *model.DeletableMedia

		isDeletableInput := c.FormValue(model.MultipartIsDeletableKey)
		if strings.ToLower(isDeletableInput) == "true" {
			isDeleteable = true

			deletable, err = s.extractDeleteRule(c)
			if err != nil {
				return s.apiRespGenerator.GenerateEchoAPIResponse(c, &stdhttp.StandardResponse{
					Success: false,
					Message: err.Error(),
					Status:  http.StatusBadRequest,
					Error:   echo.ErrBadRequest.Error(),
				}, nil)
			}
		}

		uploadInput := &model.PublicUploadInput{
			File:           file,
			IsDeletable:    isDeleteable,
			DeletableMedia: deletable,
		}

		f, err := s.publicHandler.Upload(c.Request().Context(), uploadInput)
		if err != nil {
			return s.apiRespGenerator.GenerateEchoAPIResponse(c, &stdhttp.StandardResponse{
				Success: false,
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
				Error:   echo.ErrInternalServerError.Error(),
			}, nil)
		}

		return s.apiRespGenerator.GenerateEchoAPIResponse(c, &stdhttp.StandardResponse{
			Success: true,
			Message: "success",
			Status:  http.StatusOK,
			Data:    f,
		}, nil)
	}
}

func (s *Service) handlePublicDownload() echo.HandlerFunc {
	return func(c echo.Context) error {
		filename := c.Param("filename")
		if filename == "" {
			return c.NoContent(http.StatusNotFound)
		}

		f, dec, err := s.publicHandler.Download(c.Request().Context(), filename)
		if err != nil {
			return s.apiRespGenerator.GenerateEchoAPIResponse(c, &stdhttp.StandardResponse{
				Success: false,
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
				Error:   echo.ErrInternalServerError.Error(),
			}, nil)
		}

		return c.Blob(http.StatusOK, f.ContentType, dec)
	}
}

func (s *Service) extractDeleteRule(c echo.Context) (*model.DeletableMedia, error) {
	deleteRuleInput := c.FormValue(model.MultipartDeleteRuleKey)
	deleteRule, err := model.ParseStringToDeleteRule(deleteRuleInput)
	if err != nil {
		return nil, err
	}

	switch deleteRule {
	default:
		return nil, errors.New("invalid delete rule")

	case model.ManualDelete:
		return &model.DeletableMedia{
			DeleteRule: deleteRule,
		}, nil

	case model.ScheduledDelete:
		dur := c.FormValue(model.MultipartScheduledDeleteDuration)
		if dur == "" {
			return nil, errors.New("missing scheduled delete duration")
		}

		duration, err := time.ParseDuration(dur)
		if err != nil {
			return nil, err
		}

		deleteAfter := time.Now().Add(duration)

		return &model.DeletableMedia{
			DeleteRule: deleteRule,
			Metadata: &model.DeletableRuleMetadata{
				DeleteAfter: &deleteAfter,
			},
		}, nil
	}

}
