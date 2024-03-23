package log

import "fmt"

type Level int

const (
	NoneLevel Level = iota
	ErrorLevel
	WarningLevel
	InfoLevel
	DebugLevel
)

var logLevel Level = DebugLevel

func SetGlobalLevel(level Level) {
	logLevel = level
}

// Output debug-level log message
func Debug(format string, args ...interface{}) {
	if logLevel < DebugLevel {
		return
	}

	fmt.Printf(format+"\n", args...)
}

// Output error-level log message
func Error(format string, args ...interface{}) {
	if logLevel < ErrorLevel {
		return
	}

	fmt.Printf(format+"\n", args...)
}

// Output info-level log message
func Info(format string, args ...interface{}) {
	if logLevel < InfoLevel {
		return
	}

	fmt.Printf(format+"\n", args...)
}

// Output info-level log message (in green)
func Success(format string, args ...interface{}) {
	if logLevel < InfoLevel {
		return
	}

	fmt.Printf(format+"\n", args...)
}

// Output warning-level log message
func Warning(format string, args ...interface{}) {
	if logLevel < WarningLevel {
		return
	}

	fmt.Printf(format+"\n", args...)
}
