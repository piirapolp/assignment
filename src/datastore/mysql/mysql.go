package mysql

import (
	"context"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"time"

	"assignment/logger"
	"github.com/spf13/viper"
	gormlog "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase() {
	// Init Database Config
	userName := viper.GetString("Database.Username")
	password := viper.GetString("Database.Password")
	host := viper.GetString("Database.Host")
	port := viper.GetInt("Database.Port")
	databaseName := viper.GetString("Database.DatabaseName")
	connectionTimeout := viper.GetInt("Database.ConnectionTimeout")
	maxOpen := viper.GetInt("Database.MaxConnection")
	maxIdle := viper.GetInt("Database.MinConnection")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&charset=utf8mb4&loc=Local&timeout=%ds&readTimeout=%ds&writeTimeout=%ds&multiStatements=true",
		userName, password, host, port, databaseName, connectionTimeout, connectionTimeout, connectionTimeout,
	)

	// Map log level
	dbLogLevel := logLevelFromConfig(viper.GetString("Database.LogLevel"))

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt:            true,
		Logger:                 gormlog.Default.LogMode(dbLogLevel),
	})
	if err != nil {
		logger.Logger.Errorf("Unable to connect to database: %s\n", err)
		os.Exit(1)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		logger.Logger.Errorf("unable to get underlying sql.DB: %v", err)
		os.Exit(1)
	}

	if maxOpen > 0 {
		sqlDB.SetMaxOpenConns(maxOpen)
	}
	if maxIdle > 0 {
		sqlDB.SetMaxIdleConns(maxIdle)
	}

	sqlDB.SetConnMaxLifetime(1 * time.Hour)
	sqlDB.SetConnMaxIdleTime(30 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(connectionTimeout)*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		logger.Logger.Errorf("unable to ping mysql: %v", err)
		os.Exit(1)
	}

	logger.Logger.Infof("mysql pool is started")
}

func ShutdownDatabase() {
	if DB == nil {
		logger.Logger.Infof("mysql pool is already closed")
		return
	}

	sqlDB, err := DB.DB()
	if err != nil {
		logger.Logger.Errorf("failed to get underlying sql.DB: %v", err)
		return
	}
	if err := sqlDB.Close(); err != nil {
		logger.Logger.Errorf("failed to close mysql pool: %v", err)
		return
	}

	DB = nil
	logger.Logger.Infof("mysql pool is closed")
}

func logLevelFromConfig(level string) gormlog.LogLevel {
	switch level {
	case "silent", "none":
		return gormlog.Silent
	case "error":
		return gormlog.Error
	case "warn", "warning":
		return gormlog.Warn
	case "info", "debug", "trace":
		// GORM has Info as the most verbose common level
		return gormlog.Info
	default:
		return gormlog.Info
	}
}
