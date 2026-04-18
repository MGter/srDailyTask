package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

var (
	logLevel = DEBUG
	logPath  = "./log"
	logFile  *os.File
)

func Init(level string, path string) {
	logPath = path
	switch level {
	case "debug":
		logLevel = DEBUG
	case "info":
		logLevel = INFO
	case "warning":
		logLevel = WARNING
	case "error":
		logLevel = ERROR
	default:
		logLevel = DEBUG
	}

	if err := createLogFile(); err != nil {
		log.Printf("Failed to create log file: %v", err)
	}
}

func createLogFile() error {
	dir := filepath.Join(logPath, time.Now().Format("200601"))
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	filename := filepath.Join(dir, time.Now().Format("200601")+".log")
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	if logFile != nil {
		logFile.Close()
	}
	logFile = file
	return nil
}

func formatLog(level string, file string, line int, message string) string {
	return fmt.Sprintf("[%s][%s][%s:%d] %s",
		time.Now().Format("2006-01-02 15:04:05"),
		level,
		file,
		line,
		message)
}

func writeLog(line string) {
	if logFile != nil {
		logFile.WriteString(line + "\n")
	} else {
		log.Println(line)
	}
}

func Debug(file string, line int, format string, args ...interface{}) {
	if logLevel <= DEBUG {
		msg := fmt.Sprintf(format, args...)
		writeLog(formatLog("DEBUG", file, line, msg))
	}
}

func Info(file string, line int, format string, args ...interface{}) {
	if logLevel <= INFO {
		msg := fmt.Sprintf(format, args...)
		writeLog(formatLog("INFO", file, line, msg))
	}
}

func Warning(file string, line int, format string, args ...interface{}) {
	if logLevel <= WARNING {
		msg := fmt.Sprintf(format, args...)
		writeLog(formatLog("WARNING", file, line, msg))
	}
}

func Error(file string, line int, format string, args ...interface{}) {
	if logLevel <= ERROR {
		msg := fmt.Sprintf(format, args...)
		writeLog(formatLog("ERROR", file, line, msg))
	}
}

func Close() {
	if logFile != nil {
		logFile.Close()
	}
}