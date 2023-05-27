package config

import (
	"fmt"

	"github.com/spf13/viper"
	_ "github.com/sweet-go/stdlib/config"
)

func Env() string {
	return viper.GetString("env")
}

func LogLevel() string {
	return viper.GetString("log.level")
}

func ServerPort() string {
	return fmt.Sprintf(":%s", viper.GetString("server.port"))
}

func StoragePath() string {
	return viper.GetString("server.storage_path")
}
