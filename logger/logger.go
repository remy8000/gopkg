package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger       *zap.Logger
	level        = zap.NewAtomicLevel()
	logsEncoder  zapcore.Encoder
	levelMapping = map[string]zapcore.Level{
		"debug": zapcore.DebugLevel,
		"info":  zapcore.InfoLevel,
		"warn":  zapcore.WarnLevel,
		"error": zapcore.ErrorLevel,
	}
	service string
)

func Init(srv, lvl string, printLogsInFileNotInConsole bool) {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	// store file & display console
	logsEncoder = zapcore.NewJSONEncoder(encoderCfg)

	service = srv

	// set level
	SetLevel(lvl)

	var writer zapcore.WriteSyncer
	if printLogsInFileNotInConsole {
		// Create a log file
		logFile, err := os.OpenFile("./logs/"+service+"_"+time.Now().Format("2006-01-02")+".log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			panic(err)
		}
		writer = zapcore.AddSync(logFile)
	} else {
		writer = zapcore.AddSync(zapcore.Lock(os.Stdout))
	}

	core := zapcore.NewTee(
		zapcore.NewCore(
			logsEncoder,
			writer,
			level.Level(),
		),
	)
	logger = zap.New(core)
}

// À appeler avant la fin du programme
func Close() {
    if err := logger.Sync(); err != nil && !strings.Contains(err.Error(), "inappropriate ioctl") {
        fmt.Fprintf(os.Stderr, "error syncing logger: %v\n", err)
    }
}


// SetLevel sets the log level.
func SetLevel(l string) {
	v, ok := levelMapping[l]
	if !ok {
		v = levelMapping["info"]
	}
	level.SetLevel(v)
}

// Debug wrapper.
func Debug(msg string) {
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)
	logger.Debug(msg, zap.String("service", service), zap.String("location", location))
}

// Error wrapper.
func Error(msg string) {
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)
	logger.Error(msg, zap.String("service", service), zap.String("location", location))
}

// Fatal wrapper.
func Fatal(msg, trace string) {
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)
	logger.Fatal(msg, zap.String("service", service), zap.String("location", location), zap.String("stack_trace", trace))
}

// Info wrapper.
func Info(msg string) {
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)
	logger.Info(msg, zap.String("service", service), zap.String("location", location))
}

// Panic wrapper.
func Panic(msg, trace string) {
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)
	logger.Panic(msg, zap.String("service", service), zap.String("location", location), zap.String("stack_trace", trace))
}

// Warn wrapper.
func Warn(msg string) {
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)
	logger.Warn(msg, zap.String("service", service), zap.String("location", location))
}

// Debugf wrapper.
func Debugf(t string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)
	msg := fmt.Sprintf(t, args...)
	logger.Debug(msg, zap.String("service", service), zap.String("location", location))
}

// Errorf wrapper.
func Errorf(t string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)
	msg := fmt.Sprintf(t, args...)
	logger.Error(msg, zap.String("service", service), zap.String("location", location))
}

// Fatalf wrapper.
func Fatalf(t, trace string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)
	msg := fmt.Sprintf(t, args...)
	logger.Fatal(msg, zap.String("service", service), zap.String("location", location), zap.String("stack_trace", trace))
}

// Infof wrapper.
func Infof(t string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)
	msg := fmt.Sprintf(t, args...)
	logger.Info(msg, zap.String("service", service), zap.String("location", location))
}

// Panicf wrapper.
func Panicf(t, trace string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)
	msg := fmt.Sprintf(t, args...)
	logger.Panic(msg, zap.String("service", service), zap.String("location", location), zap.String("stack_trace", trace))
}

// Warnf wrapper.
func Warnf(t string, args ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	location := fmt.Sprintf("%s:%d", file, line)
	msg := fmt.Sprintf(t, args...)
	logger.Warn(msg, zap.String("service", service), zap.String("location", location))
}
