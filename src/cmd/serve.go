package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"assignment/interface/http"
	"assignment/logger"
	"github.com/spf13/cobra"
)

func initListenOsSignal() {
	wg.Add(1)
	go func() {
		var count int
		chanOsSignal := make(chan os.Signal, 2)
		signal.Notify(chanOsSignal, syscall.SIGTERM, os.Interrupt)

		go func() {
			// Wait Os Signal
			for getSignal := range chanOsSignal {
				// Shutdown if Interrupt or SigTerm Signal is received
				if getSignal == os.Interrupt || getSignal == syscall.SIGTERM {
					count++
					// Get Twice Signal Force Exit without Waiting Close all Components
					if count == 2 {
						logger.Logger.Info("Forcefully exiting")
						os.Exit(1)
					}

					go func() {
						if enableDatabase {
							shutdownMysql()
						}
					}()

					go func() {
						if enableInterface {
							http.ShutdownHttpServer()
						}
					}()

					logger.Logger.Info("signal SIGKILL caught. shutting down")
					logger.Logger.Info("catching SIGKILL one more time will forcefully exit")

					wg.Done()
				}
			}
			close(chanOsSignal)
		}()
	}()
}

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Base Service",
	Run: func(cmd *cobra.Command, args []string) {
		// Init Component
		initComponent()

		// Init Listen OS Signal
		initListenOsSignal()

		// Init Database
		if enableDatabase {
			initMysql()
		}

		// Init Interface
		if enableInterface {
			initListenInterface()
		}

		logger.Logger.Info("service is running")

		// Waiting for Component Shut Down
		wg.Wait()

		// Flush Log
		logger.SyncLogger()
	},
}

func init() {
	rootCmd.AddCommand(ServeCmd)
}
