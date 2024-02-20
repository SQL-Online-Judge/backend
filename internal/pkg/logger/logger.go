package logger

import (
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	var (
		err      error
		logLevel = pflag.StringP("log-level", "l", "info", "set log level")
	)

	pflag.Parse()

	if *logLevel == "debug" {
		Logger, err = zap.NewDevelopment()
		Logger.Info("Log level set to debug")
	} else {
		Logger, err = zap.NewProduction()
		Logger.Info("Log level set to info")
	}

	if err != nil {
		panic(err)
	}
}
