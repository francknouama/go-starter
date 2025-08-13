package generator

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/francknouama/go-starter/internal/logger"
)

// LogLevel defines different logging levels for the generator
type LogLevel int

const (
	LogLevelSilent LogLevel = iota // No output at all
	LogLevelError                  // Errors only  
	LogLevelWarn                   // Warnings and errors
	LogLevelInfo                   // Info, warnings, and errors (default)
	LogLevelDebug                  // All logging including debug
)

// String returns the string representation of the log level
func (l LogLevel) String() string {
	switch l {
	case LogLevelSilent:
		return "silent"
	case LogLevelError:
		return "error"
	case LogLevelWarn:
		return "warn"
	case LogLevelInfo:
		return "info"
	case LogLevelDebug:
		return "debug"
	default:
		return "info"
	}
}

// GeneratorLogger provides structured logging for the generator
type GeneratorLogger struct {
	logger    logger.Logger
	level     LogLevel
	context   map[string]interface{}
	startTime time.Time
}

// NewGeneratorLogger creates a new generator logger
func NewGeneratorLogger(logLevel LogLevel) *GeneratorLogger {
	// Create a logger configuration based on level
	config := logger.DefaultConfig()
	
	// Adjust logger level based on generator level
	switch logLevel {
	case LogLevelSilent:
		config.Level = "error"
		config.Output = "/dev/null"
	case LogLevelError:
		config.Level = "error"
	case LogLevelWarn:
		config.Level = "warn"  
	case LogLevelInfo:
		config.Level = "info"
	case LogLevelDebug:
		config.Level = "debug"
	}
	
	// Use console format for better UX during generation
	config.Format = "console"
	
	baseLogger, err := logger.NewLogger(config)
	if err != nil {
		// Fallback to basic logging if logger creation fails
		baseLogger = &fallbackLogger{level: logLevel}
	}
	
	return &GeneratorLogger{
		logger:    baseLogger,
		level:     logLevel,
		context:   make(map[string]interface{}),
		startTime: time.Now(),
	}
}

// WithContext adds context to the logger
func (gl *GeneratorLogger) WithContext(ctx context.Context) *GeneratorLogger {
	return &GeneratorLogger{
		logger:    gl.logger.WithContext(ctx),
		level:     gl.level,
		context:   gl.context,
		startTime: gl.startTime,
	}
}

// WithField adds a single field to the logger context
func (gl *GeneratorLogger) WithField(key string, value interface{}) *GeneratorLogger {
	newContext := make(map[string]interface{})
	for k, v := range gl.context {
		newContext[k] = v
	}
	newContext[key] = value
	
	return &GeneratorLogger{
		logger:    gl.logger.WithFields(logger.Fields(newContext)),
		level:     gl.level,
		context:   newContext,
		startTime: gl.startTime,
	}
}

// WithFields adds multiple fields to the logger context
func (gl *GeneratorLogger) WithFields(fields map[string]interface{}) *GeneratorLogger {
	newContext := make(map[string]interface{})
	for k, v := range gl.context {
		newContext[k] = v
	}
	for k, v := range fields {
		newContext[k] = v
	}
	
	return &GeneratorLogger{
		logger:    gl.logger.WithFields(logger.Fields(newContext)),
		level:     gl.level,
		context:   newContext,
		startTime: gl.startTime,
	}
}

// Debug logs a debug message
func (gl *GeneratorLogger) Debug(msg string, args ...interface{}) {
	if gl.level >= LogLevelDebug {
		gl.logger.Debug(fmt.Sprintf(msg, args...))
	}
}

// Info logs an info message with emoji for better UX
func (gl *GeneratorLogger) Info(msg string, args ...interface{}) {
	if gl.level >= LogLevelInfo {
		gl.logger.Info(fmt.Sprintf(msg, args...))
	}
}

// Warn logs a warning message
func (gl *GeneratorLogger) Warn(msg string, args ...interface{}) {
	if gl.level >= LogLevelWarn {
		gl.logger.Warn(fmt.Sprintf(msg, args...))
	}
}

// Error logs an error message
func (gl *GeneratorLogger) Error(msg string, args ...interface{}) {
	if gl.level >= LogLevelError {
		gl.logger.Error(fmt.Sprintf(msg, args...))
	}
}

// Success logs a success message (info level with success formatting)
func (gl *GeneratorLogger) Success(msg string, args ...interface{}) {
	if gl.level >= LogLevelInfo {
		gl.logger.InfoWith(fmt.Sprintf("✓ %s", fmt.Sprintf(msg, args...)), logger.Fields{
			"type": "success",
		})
	}
}

// Progress logs progress information
func (gl *GeneratorLogger) Progress(msg string, args ...interface{}) {
	if gl.level >= LogLevelInfo {
		gl.logger.InfoWith(fmt.Sprintf("🚀 %s", fmt.Sprintf(msg, args...)), logger.Fields{
			"type": "progress",
		})
	}
}

// Warning logs a warning with emoji
func (gl *GeneratorLogger) Warning(msg string, args ...interface{}) {
	if gl.level >= LogLevelWarn {
		gl.logger.WarnWith(fmt.Sprintf("⚠️  %s", fmt.Sprintf(msg, args...)), logger.Fields{
			"type": "warning",
		})
	}
}

// ErrorWithDetails logs an error with structured details
func (gl *GeneratorLogger) ErrorWithDetails(msg string, err error, details map[string]interface{}) {
	if gl.level >= LogLevelError {
		fields := logger.Fields{
			"type":  "error",
		}
		
		if err != nil {
			fields["error"] = err.Error()
		}
		
		for k, v := range details {
			fields[k] = v
		}
		
		gl.logger.ErrorWith(fmt.Sprintf("✗ %s", msg), fields)
	}
}

// Step logs a step in the generation process
func (gl *GeneratorLogger) Step(step int, total int, msg string, args ...interface{}) {
	if gl.level >= LogLevelInfo {
		progress := float64(step) / float64(total) * 100
		gl.logger.InfoWith(fmt.Sprintf("[%.0f%%] %s", progress, fmt.Sprintf(msg, args...)), logger.Fields{
			"step":     step,
			"total":    total,
			"progress": progress,
			"type":     "step",
		})
	}
}

// Duration logs the time elapsed since logger creation
func (gl *GeneratorLogger) Duration() time.Duration {
	return time.Since(gl.startTime)
}

// Sync flushes any buffered log entries
func (gl *GeneratorLogger) Sync() error {
	return gl.logger.Sync()
}

// fallbackLogger provides basic logging when the main logger fails
type fallbackLogger struct {
	level LogLevel
}

func (fl *fallbackLogger) Debug(msg string, args ...interface{}) {
	if fl.level >= LogLevelDebug {
		fmt.Printf("[DEBUG] %s\n", fmt.Sprintf(msg, args...))
	}
}

func (fl *fallbackLogger) Info(msg string, args ...interface{}) {
	if fl.level >= LogLevelInfo {
		fmt.Printf("[INFO] %s\n", fmt.Sprintf(msg, args...))
	}
}

func (fl *fallbackLogger) Warn(msg string, args ...interface{}) {
	if fl.level >= LogLevelWarn {
		fmt.Fprintf(os.Stderr, "[WARN] %s\n", fmt.Sprintf(msg, args...))
	}
}

func (fl *fallbackLogger) Error(msg string, args ...interface{}) {
	if fl.level >= LogLevelError {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", fmt.Sprintf(msg, args...))
	}
}

func (fl *fallbackLogger) DebugWith(msg string, fields logger.Fields)              { fl.Debug("%s", msg) }
func (fl *fallbackLogger) InfoWith(msg string, fields logger.Fields)               { fl.Info("%s", msg) }
func (fl *fallbackLogger) WarnWith(msg string, fields logger.Fields)               { fl.Warn("%s", msg) }
func (fl *fallbackLogger) ErrorWith(msg string, fields logger.Fields)              { fl.Error("%s", msg) }
func (fl *fallbackLogger) WithContext(ctx context.Context) logger.Logger           { return fl }
func (fl *fallbackLogger) WithFields(fields logger.Fields) logger.Logger           { return fl }
func (fl *fallbackLogger) SetLevel(level logger.Level)                             {}
func (fl *fallbackLogger) SetOutput(w io.Writer)                                   {}
func (fl *fallbackLogger) DisableColor()                                           {}
func (fl *fallbackLogger) Sync() error                                             { return nil }

// ProgressTracker tracks progress of multi-step operations
type ProgressTracker struct {
	logger      *GeneratorLogger
	totalSteps  int
	currentStep int
	stepNames   []string
}

// NewProgressTracker creates a new progress tracker
func NewProgressTracker(logger *GeneratorLogger, stepNames []string) *ProgressTracker {
	return &ProgressTracker{
		logger:      logger,
		totalSteps:  len(stepNames),
		currentStep: 0,
		stepNames:   stepNames,
	}
}

// NextStep advances to the next step and logs progress
func (pt *ProgressTracker) NextStep() {
	if pt.currentStep < pt.totalSteps {
		pt.currentStep++
		if pt.currentStep <= len(pt.stepNames) {
			stepName := pt.stepNames[pt.currentStep-1]
			pt.logger.Step(pt.currentStep, pt.totalSteps, "%s", stepName)
		}
	}
}

// Complete logs completion of all steps
func (pt *ProgressTracker) Complete(msg string, args ...interface{}) {
	pt.logger.Success(msg, args...)
}

// Error logs an error during progress
func (pt *ProgressTracker) Error(msg string, err error) {
	pt.logger.ErrorWithDetails(msg, err, map[string]interface{}{
		"step":        pt.currentStep,
		"total_steps": pt.totalSteps,
	})
}

// LoggerOptions defines configuration options for the generator logger
type LoggerOptions struct {
	Level          LogLevel
	ShowTimestamps bool
	ShowColors     bool
	Format         string // "json", "text", "console"
	Output         string // "stdout", "stderr", file path
}

// DefaultLoggerOptions returns default logger options
func DefaultLoggerOptions() *LoggerOptions {
	return &LoggerOptions{
		Level:          LogLevelInfo,
		ShowTimestamps: false,
		ShowColors:     true,
		Format:         "console",
		Output:         "stdout",
	}
}

// SetQuietMode configures the logger for minimal output
func (opts *LoggerOptions) SetQuietMode() *LoggerOptions {
	opts.Level = LogLevelWarn
	return opts
}

// SetVerboseMode configures the logger for detailed output  
func (opts *LoggerOptions) SetVerboseMode() *LoggerOptions {
	opts.Level = LogLevelDebug
	opts.ShowTimestamps = true
	return opts
}

// SetSilentMode configures the logger for no output
func (opts *LoggerOptions) SetSilentMode() *LoggerOptions {
	opts.Level = LogLevelSilent
	return opts
}