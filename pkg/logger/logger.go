package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	levelNames = map[string]Level{
		"debug": DEBUG,
		"info":  INFO,
		"warn":  WARN,
		"error": ERROR,
	}
	levelStrings = map[Level]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
	}
)

type Logger struct {
	level  Level
	logger *log.Logger
	mu     sync.Mutex
}

var (
	defaultLogger *Logger
	once          sync.Once
)

func GetLogger() *Logger {
	once.Do(func() {
		defaultLogger = &Logger{
			level:  INFO,
			logger: log.New(os.Stdout, "", log.LstdFlags),
		}
	})
	return defaultLogger
}

func Setup(levelStr string) {
	l := GetLogger()
	l.mu.Lock()
	defer l.mu.Unlock()

	level, ok := levelNames[strings.ToLower(levelStr)]
	if !ok {
		// Default to INFO if invalid or empty
		level = DEBUG
	}
	l.level = level
}

// Helper to check if we should log
func (l *Logger) shouldLog(lvl Level) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return lvl >= l.level
}

func (l *Logger) log(lvl Level, format string, v ...interface{}) {
	if l.shouldLog(lvl) {
		msg := fmt.Sprintf(format, v...)
		prefix := fmt.Sprintf("[%s] ", levelStrings[lvl])
		l.logger.Output(3, prefix+msg) // 3 simulates calldepth to skip this wrapper
	}
}

func Debug(format string, v ...interface{}) {
	GetLogger().log(DEBUG, format, v...)
}

func Info(format string, v ...interface{}) {
	GetLogger().log(INFO, format, v...)
}

func Warn(format string, v ...interface{}) {
	GetLogger().log(WARN, format, v...)
}

func Error(format string, v ...interface{}) {
	GetLogger().log(ERROR, format, v...)
}

// For compatibility if needed, direct println style
func Debugln(v ...interface{}) {
	Debug(fmt.Sprint(v...))
}

func Infoln(v ...interface{}) {
	Info(fmt.Sprint(v...))
}

func Warnln(v ...interface{}) {
	Warn(fmt.Sprint(v...))
}

func Errorln(v ...interface{}) {
	Error(fmt.Sprint(v...))
}
