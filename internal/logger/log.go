package logger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/samber/slog-multi"
)

const (
	levelFatal = slog.Level(12)
)

var levelNames = map[slog.Leveler]string{
	levelFatal: "FATAL",
}

var LogLevelFlag string

type myLogger struct {
	Debug func(msg string, args ...any)
	Info  func(msg string, args ...any)
	Warn  func(msg string, args ...any)
	Error func(msg string, args ...any)
	Fatal func(msg string, args ...any)
}

func Get() myLogger {
	var logLevel slog.Level

	switch strings.ToUpper(LogLevelFlag) {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARNING":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	case "FATAL":
		logLevel = levelFatal
	default:
		logLevel = slog.LevelInfo
	}

	now := time.Now().UTC().Format("2006-01-02")

	logDir := "logs"
	logFile := fmt.Sprintf("server-%s.log", now)
	logPath := filepath.Join(logDir, logFile)

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
	}

	// TODO: include request id
	all := newJSONHandler(file, logLevel)
	error := newJSONHandler(os.Stderr, slog.LevelError)

	logger := *slog.New(slogmulti.Fanout(all, error))

	var ctx = context.Background()

	myLogger := myLogger{
		Debug: logger.Debug,
		Info:  logger.Info,
		Warn:  logger.Warn,
		Error: logger.Error,
		Fatal: func(msg string, args ...any) {
			logger.Log(ctx, levelFatal, msg, args...)
			os.Exit(1)
		},
	}

	return myLogger
}

func newJSONHandler(w io.Writer, level slog.Level) slog.Handler {
	return slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
		// search the custom log level name, like "FATAL"
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := levelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
	})
}
