package logging

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.SugaredLogger
var atom = zap.NewAtomicLevel()

func init() {
	encoderCfg := zap.NewDevelopmentEncoderConfig()

	logger = zap.New(zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	)).Sugar()
}

func GetLogger() *zap.SugaredLogger {
	return logger
}

func ChangeLogLevel(l zapcore.Level) {
	atom.SetLevel(l)
}
