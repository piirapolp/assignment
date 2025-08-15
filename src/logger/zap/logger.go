package zaplogger

import (
	"os"

	"assignment/logger"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger *zap.SugaredLogger
}

func getZapLevel(level string) zapcore.Level {
	switch level {
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "debug":
		return zapcore.DebugLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func NewLogger() *ZapLogger {
	logLevel := getZapLevel(viper.GetString("Log.Level"))
	logColor := viper.GetBool("Log.Color")
	logJson := viper.GetBool("Log.Json")

	var logEncoder zapcore.Encoder
	logEncoderConfig := zap.NewProductionEncoderConfig()
	logEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if logColor {
		logEncoderConfig.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	}

	if logJson {
		logEncoder = zapcore.NewJSONEncoder(logEncoderConfig)
	} else {
		logEncoder = zapcore.NewConsoleEncoder(logEncoderConfig)
	}

	core := zapcore.NewCore(
		logEncoder,
		os.Stdout,
		zap.NewAtomicLevelAt(logLevel),
	)
	zapLogger := zap.New(core, zap.AddCaller()).Sugar()

	zapLogger.Infof("init logger complete")

	return &ZapLogger{
		logger: zapLogger,
	}
}

func (l ZapLogger) GetLogger() *zap.SugaredLogger {
	return l.logger
}

func (l ZapLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l ZapLogger) Debugf(template string, args ...interface{}) {
	l.logger.Debugf(template, args...)
}

func (l ZapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

func (l ZapLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l ZapLogger) Infof(template string, args ...interface{}) {
	l.logger.Infof(template, args...)
}

func (l ZapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

func (l ZapLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l ZapLogger) Errorf(template string, args ...interface{}) {
	l.logger.Errorf(template, args...)
}

func (l ZapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

func (l ZapLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l ZapLogger) Fatalf(template string, args ...interface{}) {
	l.logger.Fatalf(template, args...)
}

func (l ZapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

func (l ZapLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l ZapLogger) Panicf(template string, args ...interface{}) {
	l.logger.Panicf(template, args...)
}

func (l ZapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.logger.Panicw(msg, keysAndValues...)
}

func (l ZapLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l ZapLogger) Warnf(template string, args ...interface{}) {
	l.logger.Warnf(template, args...)
}

func (l ZapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

func (l ZapLogger) Sync() error {
	l.logger.Infof("flush logger")
	return l.logger.Sync()
}

func (l ZapLogger) With(key string, value interface{}) logger.LoggerIface {
	return &ZapLogger{
		logger: l.logger.With(key, value),
	}
}
