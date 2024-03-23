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

// Print the text in a specified color (also adds a newline)
func colorPrintf(colorCode string, format string, args ...interface{}) {
	fmt.Printf("\033["+colorCode+"m"+format+"\033[0m\n", args...)
}

// Output debug-level log message
func Debug(format string, args ...interface{}) {
	if logLevel < DebugLevel {
		return
	}

	colorPrintf("3;37", format, args...)
}

// Output error-level log message
func Error(format string, args ...interface{}) {
	if logLevel < ErrorLevel {
		return
	}

	colorPrintf("1;31", format, args...)
}

// Output info-level log message
func Info(format string, args ...interface{}) {
	if logLevel < InfoLevel {
		return
	}

	colorPrintf("0", format, args...)
}

// Output info-level log message (in green)
func Success(format string, args ...interface{}) {
	if logLevel < InfoLevel {
		return
	}

	colorPrintf("1;32", format, args...)
}

// Output warning-level log message
func Warning(format string, args ...interface{}) {
	if logLevel < WarningLevel {
		return
	}

	colorPrintf("1;33", format+"\n", args...)
}
