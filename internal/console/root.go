package console

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sweet-go/file-server/internal/config"
	"github.com/sweet-go/stdlib/cmd"
)

func init() {
	cmd.SetupLogger(config.Env(), config.LogLevel(), "")
}

var RootCMD = cmd.CobraInitializer()

// Execute execute command
func Execute() {
	if err := RootCMD.Execute(); err != nil {
		logrus.Error(err)
		os.Exit(1)
	}
}
