// Package config provides application configuration
package config

import (
	"fmt"
	"strings"

	"github.com/hibiken/asynq"
	"github.com/spf13/viper"

	// this blank import is used to run the init function of this stdlib/config package
	_ "github.com/sweet-go/stdlib/config"
)

// Env returns application environment
func Env() string {
	return viper.GetString("env")
}

// LogLevel returns application log level
func LogLevel() string {
	return viper.GetString("log.level")
}

// ServerPort returns application server port
func ServerPort() string {
	return fmt.Sprintf(":%s", viper.GetString("server.port"))
}

// StoragePath returns application storage path
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

// WorkerBrokerHost returns worker broker host
func WorkerBrokerHost() string {
	return viper.GetString("worker.broker_host")
}

// WorkerConcurency returns worker concurrency. Default to 10
func WorkerConcurency() int {
	num := viper.GetInt("worker.concurrency")
	if num <= 0 {
		return 10
	}

	return num
}

// WorkerLogLevel returns worker log level. Default to info
func WorkerLogLevel() asynq.LogLevel {
	level := viper.GetString("worker.log_level")
	switch strings.ToUpper(level) {
	default:
		return asynq.InfoLevel
	case "DEBUG":
		return asynq.DebugLevel
	case "INFO":
		return asynq.InfoLevel
	case "WARN":
		return asynq.WarnLevel
	case "ERROR":
		return asynq.ErrorLevel
	case "FATAL":
		return asynq.FatalLevel
	}
}

// ScheduledDeleteSchedulerCronspec returns scheduled delete scheduler cronspec. Default to every 1 minute
func ScheduledDeleteSchedulerCronspec() string {
	spec := viper.GetString("worker.task.scheduled_delete.cronspec")
	if spec == "" {
		return "* * * * *"
	}

	return spec
}

// ScheduledDeleteBatchSize returns scheduled delete batch size. Default to 100
func ScheduledDeleteBatchSize() int {
	size := viper.GetInt("worker.task.scheduled_delete.batch_size")
	if size <= 0 {
		return 100
	}

	return size
}
