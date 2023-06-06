package console

import (
	"context"
	"crypto"
	"crypto/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
	"github.com/sweet-go/file-server/internal/config"
	"github.com/sweet-go/file-server/internal/delivery/httpsvc"
	"github.com/sweet-go/file-server/internal/usecase"
	"github.com/sweet-go/stdlib/encryption"
	stdhttp "github.com/sweet-go/stdlib/http"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "start the server",
	Run:   server,
}

func init() {
	RootCMD.AddCommand(serverCmd)
}

func server(cmd *cobra.Command, args []string) {
	key, err := encryption.ReadKeyFromFile("./private.pem")
	if err != nil {
		panic(err)
	}

	publicHandler := usecase.NewPublicHandler(key, config.StoragePath())
	apirespGen := stdhttp.NewStandardAPIResponseGenerator(&encryption.SignOpts{
		Random:  rand.Reader,
		PrivKey: key.PrivateKey,
		Alg:     crypto.SHA512,
		PSSOpts: nil,
	})

	HTTPServer := echo.New()

	HTTPServer.Pre(middleware.AddTrailingSlash())
	HTTPServer.Use(middleware.Recover())
	HTTPServer.Use(middleware.Logger())
	HTTPServer.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))

	publicGroup := HTTPServer.Group("public")

	httpsvc.NewService(publicGroup, publicHandler, apirespGen)

	// Start server
	go func() {
		if err := HTTPServer.Start(config.ServerPort()); err != nil && err != http.ErrServerClosed {
			HTTPServer.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := HTTPServer.Shutdown(ctx); err != nil {
		HTTPServer.Logger.Fatal(err)
	}
}
