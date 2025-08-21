package trace

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// LogLevel represents different logging levels
type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
)

// String returns the string representation of a log level
func (l LogLevel) String() string {
	switch l {
	case TRACE:
		return "TRACE"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

// Logger provides structured trace logging functionality
type Logger struct {
	enabled   bool
	minLevel  LogLevel
	output    io.Writer
	logger    *log.Logger
	startTime time.Time
}

// Global logger instance
var globalLogger *Logger

// Init initializes the trace logging system
func Init(enabled bool, filename string) error {
	globalLogger = &Logger{
		enabled:   enabled,
		minLevel:  TRACE,
		startTime: time.Now(),
	}

	if !enabled {
		globalLogger.output = io.Discard
		globalLogger.logger = log.New(io.Discard, "", 0)
		return nil
	}

	// Open trace file
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open trace file %s: %v", filename, err)
	}

	globalLogger.output = file
	globalLogger.logger = log.New(file, "", 0)

	// Log session start
	globalLogger.logf(INFO, "system", "=== CADUBI TRACE SESSION START ===")
	globalLogger.logf(INFO, "system", "Timestamp: %s", time.Now().Format("2006-01-02 15:04:05"))
	globalLogger.logf(INFO, "system", "PID: %d", os.Getpid())

	return nil
}

// Close closes the trace logging system
func Close() {
	if globalLogger != nil && globalLogger.enabled {
		globalLogger.logf(INFO, "system", "=== CADUBI TRACE SESSION END ===")
		globalLogger.logf(INFO, "system", "Session duration: %s", time.Since(globalLogger.startTime))
		
		if closer, ok := globalLogger.output.(io.Closer); ok {
			closer.Close()
		}
	}
}

// IsEnabled returns whether trace logging is enabled
func IsEnabled() bool {
	return globalLogger != nil && globalLogger.enabled
}

// logf is the internal logging function
func (l *Logger) logf(level LogLevel, category string, format string, args ...interface{}) {
	if !l.enabled || level < l.minLevel {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] [%s] [%s] %s", timestamp, level.String(), category, message)
	
	l.logger.Println(logLine)
}

// Public logging functions
func Tracef(category string, format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.logf(TRACE, category, format, args...)
	}
}

func Debugf(category string, format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.logf(DEBUG, category, format, args...)
	}
}

func Infof(category string, format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.logf(INFO, category, format, args...)
	}
}

func Warnf(category string, format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.logf(WARN, category, format, args...)
	}
}

func Errorf(category string, format string, args ...interface{}) {
	if globalLogger != nil {
		globalLogger.logf(ERROR, category, format, args...)
	}
}

// Timer provides performance timing functionality
type Timer struct {
	category  string
	operation string
	startTime time.Time
}

// StartTimer starts a performance timer
func StartTimer(category, operation string) *Timer {
	t := &Timer{
		category:  category,
		operation: operation,
		startTime: time.Now(),
	}
	Tracef(category, "Started %s", operation)
	return t
}

// End ends the timer and logs the duration
func (t *Timer) End() {
	duration := time.Since(t.startTime)
	Tracef(t.category, "Completed %s in %v", t.operation, duration)
}

// LogTerminalInfo logs terminal environment information
func LogTerminalInfo(width, height int) {
	if !IsEnabled() {
		return
	}
	
	Infof("terminal", "Terminal size: %dx%d", width, height)
	
	// Check for environment variables
	term := os.Getenv("TERM")
	colorterm := os.Getenv("COLORTERM")
	tmux := os.Getenv("TMUX")
	
	Infof("terminal", "TERM=%s, COLORTERM=%s", term, colorterm)
	if tmux != "" {
		Infof("terminal", "Running in tmux: %s", tmux)
	}
}

// LogCommand logs command execution
func LogCommand(cmdType string, details map[string]interface{}) {
	if !IsEnabled() {
		return
	}
	
	var parts []string
	for k, v := range details {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	
	if len(parts) > 0 {
		Tracef("command", "Executing %s: %s", cmdType, fmt.Sprintf("%v", parts))
	} else {
		Tracef("command", "Executing %s", cmdType)
	}
}

// LogUIEvent logs UI interactions
func LogUIEvent(eventType string, details map[string]interface{}) {
	if !IsEnabled() {
		return
	}
	
	var parts []string
	for k, v := range details {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	
	if len(parts) > 0 {
		Debugf("ui", "%s: %s", eventType, fmt.Sprintf("%v", parts))
	} else {
		Debugf("ui", "%s", eventType)
	}
}

// LogError logs error with context
func LogError(category string, err error, context map[string]interface{}) {
	if !IsEnabled() || err == nil {
		return
	}
	
	var parts []string
	for k, v := range context {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	
	if len(parts) > 0 {
		Errorf(category, "Error: %v, Context: %s", err, fmt.Sprintf("%v", parts))
	} else {
		Errorf(category, "Error: %v", err)
	}
}