package db

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/sweet-go/file-server/internal/config"
	libdb "github.com/sweet-go/stdlib/db"
	"gorm.io/gorm"

	gormLogger "gorm.io/gorm/logger"
)

var (
	PostgresDB *gorm.DB
)

func InitializePostgresConn() {
	conn, err := libdb.NewPostgresDB(config.PostgresDSN())
	if err != nil {
		logrus.Error("failed to initialize postgres connection")
		os.Exit(1)
	}

	PostgresDB = conn

	switch config.LogLevel() {
	case "error":
		PostgresDB.Logger = PostgresDB.Logger.LogMode(gormLogger.Error)
	case "warn":
		PostgresDB.Logger = PostgresDB.Logger.LogMode(gormLogger.Warn)
	default:
		PostgresDB.Logger = PostgresDB.Logger.LogMode(gormLogger.Info)
	}

	logrus.Info("Connected to Postgres Database")
}
