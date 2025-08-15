package global

import (
	"os"

	"assignment/logger"
	"github.com/spf13/viper"
)

var TimeZone string

func InitVariable() {
	TimeZone = viper.GetString("System.TimeZone")
	if TimeZone == "" {
		logger.Logger.Errorf("TimeZone variable is not config")
		os.Exit(1)
	}
}
