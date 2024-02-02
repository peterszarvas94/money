package logger

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/samber/slog-multi"
)

var logger slog.Logger

// exported functions
var Debug = logger.Debug
var Info = logger.Info
var Warn = logger.Warn
var Error = logger.Error

// add FATAL log level
var ctx = context.Background()

func Fatal(msg string, args ...any) {
	logger.Log(ctx, LevelFatal, msg, args...)
	os.Exit(1)
}

const (
	LevelFatal = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelFatal: "FATAL",
}

func init() {
	var logLevelStr string

	flag.StringVar(&logLevelStr, "log", "INFO", "-log DEBUG|INFO|WARNING|ERROR")
	flag.Parse()

	var logLevel slog.Level

	switch strings.ToUpper(logLevelStr) {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARNING":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	now := time.Now().UTC().Format("2006-01-02-15-04-05")
	logDir := "logs"
	logFile := fmt.Sprintf("server-%s.log", now)
	logPath := filepath.Join(logDir, logFile)

	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening log file:", err)
	}

	all := newJSONHandler(file, logLevel)
	error := newJSONHandler(os.Stderr, slog.LevelError)

	logger = *slog.New(slogmulti.Fanout(all, error))
}

func newJSONHandler(w io.Writer, level slog.Level) slog.Handler {
	return slog.NewJSONHandler(w, &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
    // search the custom log level name, like "FATAL"
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}

				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
	})
}
