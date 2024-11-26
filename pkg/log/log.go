package log

import (
	"fmt"
	"os"
	"sync"
	"time"
)

const (
	LevelInfo  = "INFO"
	LevelWarn  = "WARN"
	LevelError = "ERROR"
	LevelFatal = "FATAL"

	ColorReset = "\033[0m"
	ColorInfo  = "\033[32m"
	ColorWarn  = "\033[33m"
	ColorError = "\033[31m"
	ColorFatal = "\033[35m"
)

var (
	colors = map[string]string{
		LevelInfo:  ColorInfo,
		LevelWarn:  ColorWarn,
		LevelError: ColorError,
		LevelFatal: ColorFatal,
	}
	mu sync.Mutex
)

func log(level, format string, v ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	color := colors[level]
	message := fmt.Sprintf(format, v...)
	fmt.Printf("%s  %s%s%s  %s\n", timestamp, color, level, ColorReset, message)
	if level == LevelFatal {
		os.Exit(1)
	}
}

func Info(format string, v ...interface{}) {
	log(LevelInfo, format, v...)
}

func Warn(format string, v ...interface{}) {
	log(LevelWarn, format, v...)
}

func Error(format string, v ...interface{}) {
	log(LevelError, format, v...)
}

func Fatal(format string, v ...interface{}) {
	log(LevelFatal, format, v...)
}
