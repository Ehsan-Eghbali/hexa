package logutil

import (
	"dariush/config"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Cache to store logged events and prevent duplicates
var (
	loggedEvents = make(map[string]bool)
	mutex        sync.Mutex
)

// Init Logrus to format logs in JSON
func Init() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

// GenerateCorrelationID generates a unique correlation ID for tracking events
func GenerateCorrelationID() string {
	return uuid.New().String()
}

// LogRelationalStart logs the start of an event with a correlation ID
func LogRelationalStart(correlationID, event string, additionalFields map[string]interface{}) *logrus.Entry {
	if !config.Get().App.Debug {
		return nil
	}

	fields := logrus.Fields{
		"event":         event,
		"correlationID": correlationID,
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
		"status":        "started",
	}
	mergeFields(fields, additionalFields)

	entry := logrus.WithFields(fields)
	entry.Info("Event started")
	return entry
}

// LogRelationalEnd logs the end or completion of an event using the correlation ID
func LogRelationalEnd(correlationID, event string, additionalFields map[string]interface{}) *logrus.Entry {
	if !config.Get().App.Debug {
		return nil
	}
	fields := logrus.Fields{
		"event":         event,
		"correlationID": correlationID,
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
		"status":        "completed",
	}
	mergeFields(fields, additionalFields)

	entry := logrus.WithFields(fields)
	entry.Info("Event Finished")
	return entry
}

// LogError logs an error with additional context and the correlation ID
func LogError(correlationID, event string, err error, additionalFields map[string]interface{}) {
	fields := logrus.Fields{
		"event":         event,
		"correlationID": correlationID,
		"timestamp":     time.Now().UTC().Format(time.RFC3339),
		"error":         err.Error(),
		"status":        "error",
	}
	mergeFields(fields, additionalFields)

	logrus.WithFields(fields).Error("Error occurred")
}

// LogOnce logs an event only if it has not been logged before to prevent duplicates
func LogOnce(event string, err error, additionalFields map[string]interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	// Check if the event has already been logged using correlationID + event as the key
	logKey := event
	if loggedEvents[logKey] {
		return // Do not log again if it already exists
	}

	fields := logrus.Fields{
		"err":       err.Error(),
		"event":     event,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	mergeFields(fields, additionalFields)

	logrus.WithFields(fields).Info("Logged once event")
	loggedEvents[logKey] = true
}

// LogSuccess logs an event only if it has not been logged before to prevent duplicates
func LogSuccess(event string, additionalFields map[string]interface{}) {
	mutex.Lock()
	defer mutex.Unlock()

	// Check if the event has already been logged using correlationID + event as the key
	logKey := event
	if loggedEvents[logKey] {
		return // Do not log again if it already exists
	}

	fields := logrus.Fields{
		"event":     event,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	mergeFields(fields, additionalFields)

	logrus.WithFields(fields).Info("Logged success event")
	loggedEvents[logKey] = true
}

// Helper to merge additional fields into the log fields
func mergeFields(baseFields logrus.Fields, additionalFields map[string]interface{}) {
	for k, v := range additionalFields {
		baseFields[k] = v
	}
}
