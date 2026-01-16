package logger

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestLoggerFiltering(t *testing.T) {
	// Capture output
	var buf bytes.Buffer
	l := GetLogger()
	// Save original state
	l.mu.Lock()
	origLogger := l.logger
	origLevel := l.level
	l.logger = log.New(&buf, "", 0)
	l.level = INFO
	l.mu.Unlock()

	// Restore after test
	defer func() {
		l.mu.Lock()
		l.logger = origLogger
		l.level = origLevel
		l.mu.Unlock()
	}()

	// reset buffer
	buf.Reset()
	Debug("This is debug")
	if buf.Len() > 0 {
		t.Errorf("Expected no output for DEBUG when level is INFO, got %s", buf.String())
	}

	buf.Reset()
	Info("This is info")
	if !strings.Contains(buf.String(), "[INFO] This is info") {
		t.Errorf("Expected INFO log, got %s", buf.String())
	}

	buf.Reset()
	Warn("This is warn")
	if !strings.Contains(buf.String(), "[WARN] This is warn") {
		t.Errorf("Expected WARN log, got %s", buf.String())
	}
}

func TestSetup(t *testing.T) {
	Setup("debug")
	if GetLogger().level != DEBUG {
		t.Errorf("Expected level DEBUG, got %v", GetLogger().level)
	}

	Setup("invalid")
	if GetLogger().level != INFO {
		t.Errorf("Expected default INFO, got %v", GetLogger().level)
	}
}
