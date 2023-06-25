package httpsvc

import (
	"net/http"

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

		f, err := s.publicHandler.Upload(c.Request().Context(), file)
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
