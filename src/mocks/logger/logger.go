package fake_logger

import (
	"assignment/logger"
)

type FakeLogger struct {
}

func NewLogger() *FakeLogger {
	return &FakeLogger{}
}

func (l FakeLogger) Debug(args ...interface{}) {
}

func (l FakeLogger) Debugf(template string, args ...interface{}) {
}

func (l FakeLogger) Debugw(msg string, keysAndValues ...interface{}) {
}

func (l FakeLogger) Info(args ...interface{}) {
}

func (l FakeLogger) Infof(template string, args ...interface{}) {
}

func (l FakeLogger) Infow(msg string, keysAndValues ...interface{}) {
}

func (l FakeLogger) Error(args ...interface{}) {
}

func (l FakeLogger) Errorf(template string, args ...interface{}) {
}

func (l FakeLogger) Errorw(msg string, keysAndValues ...interface{}) {
}

func (l FakeLogger) Fatal(args ...interface{}) {
}

func (l FakeLogger) Fatalf(template string, args ...interface{}) {
}

func (l FakeLogger) Fatalw(msg string, keysAndValues ...interface{}) {
}

func (l FakeLogger) Panic(args ...interface{}) {
}

func (l FakeLogger) Panicf(template string, args ...interface{}) {
}

func (l FakeLogger) Panicw(msg string, keysAndValues ...interface{}) {
}

func (l FakeLogger) Warn(args ...interface{}) {
}

func (l FakeLogger) Warnf(template string, args ...interface{}) {
}

func (l FakeLogger) Warnw(msg string, keysAndValues ...interface{}) {
}

func (l FakeLogger) Sync() error {
	return nil
}

func (l FakeLogger) With(key string, value interface{}) logger.LoggerIface {
	return l
}
