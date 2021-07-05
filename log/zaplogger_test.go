package log

import (
	"alertmanager_notifier/config"
	"go.uber.org/zap"
	"testing"
	"time"
)

func TestNewZapLogger(t *testing.T) {
	runConf := config.NewRunConf("debug", "dev")
	logger := NewZapLogger(ConfigZap("", runConf))

	t.Run("simple_test", func(t *testing.T) {
		logger.Info("test Info",
			zap.String("a", "aaa"),
			zap.Int("b", 3),
			zap.Duration("c", time.Second),
		)
	})

}

func TestNewZapSugarLogger(t *testing.T) {
	runConf := config.NewRunConf("debug", "dev")
	logger := NewZapSugarLogger(ConfigZap("", runConf))

	t.Run("simple_sugar_test", func(t *testing.T) {
		logger.Debug("test", "test Debug")
		logger.Info("test", "test Info")
		logger.Warn("test", "test Warn")
		logger.Error("test","test Error")
		logger.Panic("test", "test Panic")
	})
}

