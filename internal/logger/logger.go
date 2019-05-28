/*Package logger ...*/
/*
This package implements the logger which helps us to view the application logs
Application logs can be redirected to a standard ouptut or to a log file
This package supports three levels of logging DEBUG,INFO,ERROR
Debug logs will be written only if debugEnabled is set in config
This will not perform any error handling it just adds the useful information in the logs
*/
package logger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

//The Logger type is the skeleton and provided the base methods to be implemented
type Logger interface {
	Error(msg string, ctx map[string]interface{})
	Info(msg string, ctx map[string]interface{})
	Debug(msg string, ctx map[string]interface{})
}

//The logFile type implements the Logger interface and abstracts the logger and debug flag
type logFile struct {
	logger         *log.Logger
	isDebugEnabled bool
}

//New function creates the basic logger will be part of application type
func New(w io.Writer, isDebug bool) Logger {
	return &logFile{
		logger:         log.New(w, logPid, logFlags),
		isDebugEnabled: isDebug,
	}
}

const logFlags = log.Ldate | log.Ltime | log.Lmicroseconds //timestamp
var (
	logPid = fmt.Sprintf("(%d) ", os.Getpid()) //process ID
)

//This function converts this to a JSON format and writes the information
//to the writer
func (l *logFile) write(level, msg string, ctx map[string]interface{}) {
	js, err := json.Marshal(struct {
		Level   string                 `json:"level"`
		Event   string                 `json:"msg"`
		Context map[string]interface{} `json:"context"`
	}{
		level,
		msg,
		ctx,
	})
	if nil == err {
		l.logger.Printf(string(js))
	} else {
		l.logger.Printf("unable to marshal log entry: %v", err)
	}
}

//write the Error information to the log
func (l *logFile) Error(msg string, ctx map[string]interface{}) {
	l.write("error", msg, ctx)
}

//write the general information to the log
func (l *logFile) Info(msg string, ctx map[string]interface{}) {
	l.write("info", msg, ctx)
}

//write the debug information to the log, this is based on the debug flag
//if the debug flag is enabled , then write the debug information to the log
func (l *logFile) Debug(msg string, ctx map[string]interface{}) {
	if l.isDebugEnabled {
		l.write("debug", msg, ctx)
	}
}
