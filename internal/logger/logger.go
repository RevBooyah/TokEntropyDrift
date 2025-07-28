package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

// Logger wraps logrus.Logger with additional functionality
type Logger struct {
	*logrus.Logger
}

// New creates a new logger with the specified configuration
func New(level, format, file string) (*Logger, error) {
	logger := logrus.New()
	
	// Set log level
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(logLevel)
	
	// Set log format
	switch format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	default:
		logger.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
		})
	}
	
	// Set output
	if file != "" {
		file, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		logger.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
		logger.SetOutput(os.Stdout)
	}
	
	return &Logger{Logger: logger}, nil
}

// WithField adds a field to the logger
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields adds multiple fields to the logger
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// WithError adds an error field to the logger
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}

// LogAnalysisStart logs the start of an analysis run
func (l *Logger) LogAnalysisStart(inputFile string, tokenizers []string) {
	l.WithFields(logrus.Fields{
		"event":      "analysis_start",
		"input_file": inputFile,
		"tokenizers": tokenizers,
	}).Info("Starting tokenization analysis")
}

// LogAnalysisComplete logs the completion of an analysis run
func (l *Logger) LogAnalysisComplete(stats map[string]interface{}) {
	l.WithFields(logrus.Fields{
		"event": "analysis_complete",
		"stats": stats,
	}).Info("Tokenization analysis completed")
}

// LogTokenizerStart logs the start of tokenizer processing
func (l *Logger) LogTokenizerStart(tokenizerName string, inputFile string) {
	l.WithFields(logrus.Fields{
		"event":          "tokenizer_start",
		"tokenizer_name": tokenizerName,
		"input_file":     inputFile,
	}).Info("Starting tokenizer processing")
}

// LogTokenizerComplete logs the completion of tokenizer processing
func (l *Logger) LogTokenizerComplete(tokenizerName string, tokenCount int, duration float64) {
	l.WithFields(logrus.Fields{
		"event":          "tokenizer_complete",
		"tokenizer_name": tokenizerName,
		"token_count":    tokenCount,
		"duration_ms":    duration,
	}).Info("Tokenizer processing completed")
}

// LogMetricCalculation logs metric calculation events
func (l *Logger) LogMetricCalculation(metricName string, tokenizerName string, value float64) {
	l.WithFields(logrus.Fields{
		"event":          "metric_calculation",
		"metric_name":    metricName,
		"tokenizer_name": tokenizerName,
		"value":          value,
	}).Debug("Metric calculated")
}

// LogVisualizationGenerated logs visualization generation events
func (l *Logger) LogVisualizationGenerated(vizType string, outputPath string) {
	l.WithFields(logrus.Fields{
		"event":       "visualization_generated",
		"viz_type":    vizType,
		"output_path": outputPath,
	}).Info("Visualization generated")
}

// LogError logs error events with context
func (l *Logger) LogError(event string, err error, context map[string]interface{}) {
	fields := logrus.Fields{
		"event": event,
		"error": err.Error(),
	}
	for k, v := range context {
		fields[k] = v
	}
	l.WithFields(fields).Error("Error occurred")
}

// LogWarning logs warning events with context
func (l *Logger) LogWarning(event string, message string, context map[string]interface{}) {
	fields := logrus.Fields{
		"event":   event,
		"message": message,
	}
	for k, v := range context {
		fields[k] = v
	}
	l.WithFields(fields).Warn("Warning occurred")
} 