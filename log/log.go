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

func SetGlobalLevel(level Level) {
	fmt.Printf("TODO: Set log level to %d", level)
}

// Output debug-level log message
func Debug(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Output error-level log message
func Error(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Output info-level log message
func Info(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Output info-level log message (in green)
func Success(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// Output warning-level log message
func Warning(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}
