package logger

import (
	"bytes"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	var buf bytes.Buffer
	logger := New(DEBUG, &buf)

	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()
	
	if !strings.Contains(output, "DEBUG") {
		t.Error("Expected DEBUG log level in output")
	}
	if !strings.Contains(output, "INFO") {
		t.Error("Expected INFO log level in output")
	}
	if !strings.Contains(output, "WARN") {
		t.Error("Expected WARN log level in output")
	}
	if !strings.Contains(output, "ERROR") {
		t.Error("Expected ERROR log level in output")
	}
	
	if !strings.Contains(output, "debug message") {
		t.Error("Expected debug message in output")
	}
	if !strings.Contains(output, "info message") {
		t.Error("Expected info message in output")
	}
	if !strings.Contains(output, "warn message") {
		t.Error("Expected warn message in output")
	}
	if !strings.Contains(output, "error message") {
		t.Error("Expected error message in output")
	}
}

func TestLoggerLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := New(INFO, &buf)

	logger.Debug("debug message")
	logger.Info("info message")

	output := buf.String()
	
	if strings.Contains(output, "debug message") {
		t.Error("DEBUG message should be filtered out with INFO level")
	}
	if !strings.Contains(output, "info message") {
		t.Error("INFO message should be included with INFO level")
	}
}

func TestLoggerFormatting(t *testing.T) {
	var buf bytes.Buffer
	logger := New(INFO, &buf)

	logger.Info("formatted %s with %d args", "message", 2)

	output := buf.String()
	if !strings.Contains(output, "formatted message with 2 args") {
		t.Error("Expected formatted message in output")
	}
}

func TestSetVerbose(t *testing.T) {
	var buf bytes.Buffer
	oldLogger := defaultLogger
	defaultLogger = New(INFO, &buf)
	defer func() { defaultLogger = oldLogger }()

	SetVerbose(true)
	Debug("test debug")
	
	output := buf.String()
	if !strings.Contains(output, "test debug") {
		t.Error("Expected debug message after SetVerbose(true)")
	}

	buf.Reset()
	SetVerbose(false)
	Debug("test debug 2")
	
	output = buf.String()
	if strings.Contains(output, "test debug 2") {
		t.Error("Debug message should be filtered after SetVerbose(false)")
	}
}

func TestLevelString(t *testing.T) {
	tests := []struct {
		level    Level
		expected string
	}{
		{DEBUG, "DEBUG"},
		{INFO, "INFO"},
		{WARN, "WARN"},
		{ERROR, "ERROR"},
	}

	for _, test := range tests {
		if test.level.String() != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, test.level.String())
		}
	}
}
