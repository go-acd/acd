//go:generate stringer -type=Level

// Package log provides a common logging package for ACD with logging level.
package log

import (
	"fmt"
	stdLog "log"
)

// Level is a custom type representing a log level.
type Level uint8

const (
	// DisableLogLevel disables logging completely.
	DisableLogLevel Level = iota

	// FatalLevel represents a fatal message.
	FatalLevel

	// ErrorLevel represents an error message.
	ErrorLevel

	// InfoLevel represents an info message.
	InfoLevel

	// DebugLevel represents a debug message.
	DebugLevel
)

var (
	// Level defines the log level. Default: Error
	level = ErrorLevel

	levelPrefix = map[Level]string{
		DisableLogLevel: "",
		FatalLevel:      "[FATAL] ",
		ErrorLevel:      "[ERROR] ",
		InfoLevel:       "[INFO] ",
		DebugLevel:      "[DEBUG] ",
	}
)

// Levels returns a string of all possible levels
func Levels() string {
	return fmt.Sprintf("%d:%s, %d:%s, %d:%s, %d:%s, %d:%s",
		DisableLogLevel, DisableLogLevel,
		FatalLevel, FatalLevel,
		ErrorLevel, ErrorLevel,
		InfoLevel, InfoLevel,
		DebugLevel, DebugLevel)
}

// SetLevel sets the log level to l.
func SetLevel(l Level) {
	level = l
}

// GetLevel sets the log level to l.
func GetLevel() Level {
	return level
}

// Printf calls Printf only if the level is equal or lower than the set level.
// If the level is FatalLevel, it will call Fatalf regardless...
func Printf(l Level, format string, v ...interface{}) {
	if l == FatalLevel {
		stdLog.Fatalf(format, v...)
		return
	}

	if l <= level {
		defer stdLog.SetPrefix(stdLog.Prefix())
		stdLog.SetPrefix(levelPrefix[l])
		stdLog.Printf(format, v...)
	}
}

// Print calls Print only if the level is equal or lower than the set level.
// If the level is FatalLevel, it will call Fatal regardless...
func Print(l Level, v ...interface{}) {
	if l == FatalLevel {
		stdLog.Fatal(v...)
		return
	}

	if l <= level {
		defer stdLog.SetPrefix(stdLog.Prefix())
		stdLog.SetPrefix(levelPrefix[l])
		stdLog.Print(v...)
	}
}

// Fatalf wraps Printf
func Fatalf(format string, v ...interface{}) { Printf(FatalLevel, format, v...) }

// Errorf wraps Printf
func Errorf(format string, v ...interface{}) { Printf(ErrorLevel, format, v...) }

// Infof wraps Printf
func Infof(format string, v ...interface{}) { Printf(InfoLevel, format, v...) }

// Debugf wraps Printf
func Debugf(format string, v ...interface{}) { Printf(DebugLevel, format, v...) }

// Fatal wraps Print
func Fatal(v ...interface{}) { Print(FatalLevel, v...) }

// Error wraps Print
func Error(v ...interface{}) { Print(ErrorLevel, v...) }

// Info wraps Print
func Info(v ...interface{}) { Print(InfoLevel, v...) }

// Debug wraps Print
func Debug(v ...interface{}) { Print(DebugLevel, v...) }
