package log

import (
	golog "log"
	"os"
)

var logger *golog.Logger

// DebugEnabled if true debug message are printed to the console
var DebugEnabled bool

func init() {
	logger = golog.New(os.Stdout, "", 0)
}

// Debugf log a debug message, if DebugEnabled is true
func Debugf(format string, v ...interface{}) {
	if !DebugEnabled {
		return
	}

	logger.Printf(format, v...)
}

// Debugln log a debug message, if DebugEnabled is true
func Debugln(v ...interface{}) {
	if !DebugEnabled {
		return
	}

	logger.Println(v...)
}

// Errorf log an error message
func Errorf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

// Infof log a Info message
func Infof(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

// Infoln log a Info message
func Infoln(v ...interface{}) {
	logger.Println(v...)
}

// Fatalf log a message and terminate the program
func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

// Fatalln log a message and terminate the program
func Fatalln(v ...interface{}) {
	logger.Fatalln(v...)
}
