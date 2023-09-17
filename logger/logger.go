package logger

import (
	"go-micro.dev/v4/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"

	logger_zap "github.com/go-micro/plugins/v4/logger/zap"
)

func init() {
	enc := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	w := zapcore.AddSync(os.Stdout)
	core := zapcore.NewCore(enc, w, zapcore.ErrorLevel)
	zapLogger := zap.New(core)
	logger.DefaultLogger, _ = logger_zap.NewLogger(logger.WithLevel(logger.InfoLevel), logger_zap.WithLogger(zapLogger))
}

func Print(args ...interface{}) {
	Info(args)
}

func Printf(template string, args ...interface{}) {
	Infof(template, args)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infof(template string, args ...interface{}) {
	logger.Infof(template, args...)
}

func Trace(args ...interface{}) {
	logger.Trace(args...)
}

func Tracef(template string, args ...interface{}) {
	logger.Tracef(template, args...)
}

func Debug(args ...interface{}) {
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	logger.Debugf(template, args...)
}

func Warn(args ...interface{}) {
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	logger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	logger.Errorf(template, args...)
}

func Fatal(args ...interface{}) {
	logger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	logger.Fatalf(template, args...)
}
