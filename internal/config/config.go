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

// PostgresDSN returns postgres DSN
func PostgresDSN() string {
	host := viper.GetString("postgres.host")
	db := viper.GetString("postgres.db")
	user := viper.GetString("postgres.user")
	pw := viper.GetString("postgres.pw")
	port := viper.GetString("postgres.port")
	sslMode := viper.GetString("postgres.ssl_mode")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s", host, user, pw, db, port, sslMode)
}
