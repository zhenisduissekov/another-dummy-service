package log

import (
	"log"
)

var Logger = log.New(log.Writer(), "lo", log.LstdFlags)

// Info prints an informational message using the global logger.
func Info(v ...any) {
	Logger.Println(v...)
}

func Infof(format string, v ...any) {
	Logger.Printf(format, v...)
}

// Error prints an error message using the global logger.
func Error(v ...any) {
	Logger.Println(v...)
}

func Errorf(format string, v ...any) {
	Logger.Printf(format, v...)
}

// Debug prints a debug message using the global logger.
func Debug(v ...any) {
	Logger.Println(v...)
}

func Fatalf(format string, v ...any) {
	Logger.Fatalf(format, v...)
}
