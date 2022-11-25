package zap

import (
	"fmt"
	"path/filepath"

	"github.com/mr-chelyshkin/go-gridder/internal/sys"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	defaultLogLevel = "INFO"
	commonAttrKey   = "attrs"
)

// Logger object with methods.
type logger struct {
	l *zap.Logger
}

// NewLoggerStdout return zap logger object.
func NewLoggerStdout(cfg config) (*logger, error) {
	if cfg == nil {
		return nil, fmt.Errorf("incoming config object is nil")
	}

	z := zap.NewProductionConfig()
	z.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	z.DisableCaller = true

	if cfg.GetLogLevel() == "" {
		cfg.SetLogLevel(defaultLogLevel)
	}
	if err := z.Level.UnmarshalText([]byte(cfg.GetLogLevel())); err != nil {
		return nil, fmt.Errorf("failed to unmarshal log level: %s", err)
	}

	l, err := z.Build()
	if err != nil {
		return nil, err
	}
	return &logger{l: l}, nil
}

// NewLoggerFile return zap logger object.
func NewLoggerFile(cfg config) (*logger, error) {
	if cfg == nil {
		return nil, fmt.Errorf("incoming config object is nil")
	}

	z := zap.NewProductionConfig()
	z.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	z.DisableCaller = true

	if cfg.GetLogFileMaxSize() == 0 {
		cfg.SetLogFileMaxSize(100)
	}
	if cfg.GetLogFileMaxAge() == 0 {
		cfg.SetLogFileMaxAge(30)
	}
	if cfg.GetLogFileMaxBackups() == 0 {
		cfg.SetLogFileMaxBackups(3)
	}
	if cfg.GetLogLevel() == "" {
		cfg.SetLogLevel(defaultLogLevel)
	}

	if cfg.GetLogFilePath() == "" {
		return nil, fmt.Errorf("failed to initialize logger to file, path not set in config")
	}
	if err := z.Level.UnmarshalText([]byte(cfg.GetLogLevel())); err != nil {
		return nil, fmt.Errorf("failed to unmarshal log level: %s", err)
	}
	if err := sys.DirCreate(filepath.Dir(cfg.GetLogFilePath())); err != nil {
		return nil, fmt.Errorf("log dir not found, can't create, got %s", err)
	}
	if _, err := sys.DirIsWritable(filepath.Dir(cfg.GetLogFilePath())); err != nil {
		return nil, err
	}

	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   cfg.GetLogFilePath(),
		MaxAge:     cfg.GetLogFileMaxAge(),
		MaxSize:    cfg.GetLogFileMaxSize(),
		MaxBackups: cfg.GetLogFileMaxBackups(),
	})
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(z.EncoderConfig),
		w,
		z.Level,
	)
	return &logger{l: zap.New(core)}, nil
}

// Debugf logs a message at level Debug on the ZapLogger with common attrs.
func (l *logger) Debugf(msg string, i interface{}) {
	l.l.Debug(msg, zap.Any(commonAttrKey, i))
}

// Infof logs a message at level Info on the ZapLogger with common attrs.
func (l *logger) Infof(msg string, i interface{}) {
	l.l.Info(msg, zap.Any(commonAttrKey, i))
}

// Warnf logs a message at level Warn on the ZapLogger with common attrs.
func (l *logger) Warnf(msg string, i interface{}) {
	l.l.Warn(msg, zap.Any(commonAttrKey, i))
}

// Errorf logs a message at level Error on the ZapLogger with common attrs.
func (l *logger) Errorf(msg string, i interface{}) {
	l.l.Error(msg, zap.Any(commonAttrKey, i))
}

// Fatalf logs a message at level Fatal on the ZapLogger with common attrs.
func (l *logger) Fatalf(msg string, i interface{}) {
	l.l.Fatal(msg, zap.Any(commonAttrKey, i))
}

// Panicf logs a message at level Panic on the ZapLogger with common attrs.
func (l *logger) Panicf(msg string, i interface{}) {
	l.l.Panic(msg, zap.Any(commonAttrKey, i))
}

// Printf logs a message at level Info on the ZapLogger.
func (l *logger) Printf(format string, args ...interface{}) {
	l.l.Info(fmt.Sprintf(format, args...))
}
