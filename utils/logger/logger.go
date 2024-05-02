package logger

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

// LogLevel represents different log levels
type LogLevel int

const (
	// LogDebug represents debug log level
	logDebug LogLevel = iota
	// LogInfo represents info log level
	logInfo
	// LogWarning represents warning log level
	logWarning
	// LogError represents error log level
	logError
)

// logLevelToString maps LogLevel to string representation
var logLevelToString = map[LogLevel]string{
	logDebug:   "DEBUG",
	logInfo:    "INFO",
	logWarning: "WARNING",
	logError:   "ERROR",
}

// log prints log message with log level, time, caller function, and line number
func log(level LogLevel, message string) {
	logType := logLevelToString[level]
	now := time.Now().Format("2006-01-02 15:04:05")

	// Get caller function and line number
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	} else {
		file = file[strings.LastIndex(file, "/")+1:] // Get only the file name
	}

	logMessage := fmt.Sprintf("[%s] [%s] %s:%d - %s\n", now, logType, file, line, message)
	fmt.Println(logMessage)
}

func LogDebug(message string) {
	log(0, message)
}
func LogInfo(message string) {
	log(1, message)
}
func LogWarning(message string) {
	log(2, message)
}
func LogError(message string) {
	log(3, message)
}
