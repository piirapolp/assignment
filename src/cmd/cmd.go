package cmd

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/shopspring/decimal"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"assignment/datastore/mysql"
	"assignment/global"
	"assignment/interface/http"
	"assignment/logger"
	zaplogger "assignment/logger/zap"
)

var wg sync.WaitGroup
var configFile string
var enableDatabase bool
var enableInterface bool

var rootCmd = &cobra.Command{
	Use:   global.BASE_SERVICE_NAME,
	Short: global.BASE_SERVICE_SHORT_NAME,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func initTimezone() {
	loc, err := time.LoadLocation(global.TimeZone)
	if err != nil {
		logger.Logger.Error("unable to set timezone because: %s", err)
		os.Exit(0)
	}
	time.Local = loc
}

func initComponent() {
	// Init Logger
	logger.Logger = zaplogger.NewLogger()

	// Init Global Variable
	global.InitVariable()

	// Init Timezone
	initTimezone()
}

func initConfig() {
	// Init viper
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		viper.AddConfigPath("./config")
		viper.SetConfigName("config")
	}

	// Read Config
	if err := viper.ReadInConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to read config: %v\n", err)
		os.Exit(1)
	}

	// Set Component Enable Flag
	enableInterface = viper.GetBool("Interface.Enable")
	enableDatabase = viper.GetBool("Database.Enable")

	// Parse Json Decimal to Float
	decimal.MarshalJSONWithoutQuotes = true
}

func initMysql() {
	logger.Logger.Info("initializing mysql")
	wg.Add(1)
	mysql.InitDatabase()
}

func shutdownMysql() {
	logger.Logger.Info("shutting down mysql")
	mysql.ShutdownDatabase()
	wg.Done()
}

func initListenInterface() {
	wg.Add(1)
	go func() {
		http.InitHttpServer()
		wg.Done()
	}()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "./config/config.yaml", "config file")
}
