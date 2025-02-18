package plugin

import (
	"github.com/hashicorp/go-hclog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type LogAdapter struct {
	logger hclog.Logger
}

func NewLogAdapter(logger hclog.Logger) *zap.Logger {
	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{})
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(logWriter{logger: logger}),
		zapcore.DebugLevel,
	)

	return zap.New(core)
}

type logWriter struct {
	logger hclog.Logger
}

func (w logWriter) Write(p []byte) (n int, err error) {
	w.logger.Info(string(p))
	return len(p), nil
}

func (w logWriter) Sync() error {
	return nil
}
