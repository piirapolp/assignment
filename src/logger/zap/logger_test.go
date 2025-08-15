package zaplogger_test

import (
	"testing"

	zaplogger "assignment/logger/zap"
)

func BenchmarkZapLogger(b *testing.B) {
	logger := zaplogger.NewLogger()

	b.Run("Debug", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			logger.Debug("Debug")
		}
	})
}
