package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type LoglevelEnum int

const (
	INFO LoglevelEnum = iota
	WARNING
	ERROR
	FATAL
)

const maxLoglevel = int(FATAL)

// Returns the string representation of the loglevel
func llToString(ll LoglevelEnum) string {
	switch ll {
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "INFO"
	}
}

// Returns the Loglevel from the string representation
func llToEnum(ll string) LoglevelEnum {
	switch ll {
	case "INFO":
		return INFO
	case "WARNING":
		return WARNING
	case "ERROR":
		return ERROR
	case "FATAL":
		return FATAL
	default:
		return INFO
	}
}

// Logs into log.txt in the following format:
// LEVEL location YYYY-MM-DD HH:MM:SS - message
func Log(messageLoglevelEnum LoglevelEnum, location string, message string) {
	fileName := "internal/logger/log.txt"
	now := time.Now().Format("2006-01-02 15:04:05")

	file, fileErr := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if fileErr != nil {
		logger := log.New(os.Stdout, "", 0)
		logger.Printf("ERROR log/file %s - %s", now, fileErr.Error())
	}
	defer file.Close()

	messageLoglevelStr := llToString(messageLoglevelEnum)

	formattedMessage := fmt.Sprintf("%s %s %s - %s", messageLoglevelStr, location, now, message)
	logger := log.New(file, "", 0)

	loglevelEnum := GetLogLevel()

	if messageLoglevelEnum < loglevelEnum {
		return
	}

	if messageLoglevelEnum == FATAL {
		logger.Fatal(formattedMessage)
	} else {
		logger.Println(formattedMessage)
	}
}

var Loglevel string

// Returns the loggin level from the -l flag, defaulting to INFO
func GetLogLevel() LoglevelEnum {
	upperLoglevel := strings.ToUpper(Loglevel)
	return llToEnum(upperLoglevel)
}
