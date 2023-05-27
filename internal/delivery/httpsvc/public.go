package httpsvc

import (
	"net/http"

	"github.com/labstack/echo/v4"
	stdhttp "github.com/sweet-go/stdlib/http"
)

func (s *Service) handlePublicUpload() echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			resp, err := s.apiRespGenerator.GenerateAPIResponse(&stdhttp.StandardResponse{
				Success: false,
				Message: "missing required file",
				Status:  http.StatusBadRequest,
				Error:   echo.ErrBadRequest.Error(),
			}, nil)

			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}

			return c.JSON(http.StatusBadRequest, resp)
		}

		f, err := s.publicHandler.Upload(c.Request().Context(), file)
		if err != nil {
			resp, err := s.apiRespGenerator.GenerateAPIResponse(&stdhttp.StandardResponse{
				Success: false,
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
				Error:   echo.ErrInternalServerError.Error(),
			}, nil)

			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}

			return c.JSON(http.StatusInternalServerError, resp)
		}

		return c.JSON(http.StatusOK, f)
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
			resp, err := s.apiRespGenerator.GenerateAPIResponse(&stdhttp.StandardResponse{
				Success: false,
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
				Error:   echo.ErrInternalServerError.Error(),
			}, nil)

			if err != nil {
				return c.NoContent(http.StatusInternalServerError)
			}

			return c.JSON(http.StatusInternalServerError, resp)
		}
		return c.Blob(http.StatusOK, f.ContentType, dec)
	}
}
