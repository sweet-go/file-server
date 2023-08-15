package httpsvc

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sweet-go/file-server/model"
	stdhttp "github.com/sweet-go/stdlib/http"
)

func (s *Service) handleDeleteMedia() echo.HandlerFunc {
	return func(c echo.Context) error {
		input := &model.DeleteMediaInput{}
		if err := c.Bind(input); err != nil {
			return s.apiRespGenerator.GenerateEchoAPIResponse(c, &stdhttp.StandardResponse{
				Success: false,
				Message: "invalid input",
				Status:  http.StatusBadRequest,
				Error:   err.Error(),
			}, nil)
		}

		if err := s.deletableMediaUsecase.DeleteMedia(c.Request().Context(), input); err != nil {
			return s.apiRespGenerator.GenerateEchoAPIResponse(c, &stdhttp.StandardResponse{
				Success: false,
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
				Error:   err.Error(),
			}, nil)
		}

		return s.apiRespGenerator.GenerateEchoAPIResponse(c, &stdhttp.StandardResponse{
			Success: true,
			Message: "success",
			Status:  http.StatusOK,
		}, nil)
	}
}
