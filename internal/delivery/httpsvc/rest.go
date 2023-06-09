package httpsvc

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/sweet-go/file-server/model"
	stdhttp "github.com/sweet-go/stdlib/http"
)

type Service struct {
	publicGroup      *echo.Group
	publicHandler    model.PublicHandler
	apiRespGenerator stdhttp.APIResponseGenerator
}

func NewService(publicGroup *echo.Group, publicHandler model.PublicHandler, apiRespGenerator stdhttp.APIResponseGenerator) {
	s := &Service{
		publicGroup,
		publicHandler,
		apiRespGenerator,
	}

	s.initPublicRoutes()
}

func (s *Service) initPublicRoutes() {
	s.publicGroup.GET("/ping/", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})

	s.publicGroup.POST("/upload/", s.handlePublicUpload())
	s.publicGroup.GET("/download/:filename/", s.handlePublicDownload())
}
