package main

import (
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
)

func isVerbose(m *pubsub.Message) bool {
	if m == nil {
		return *oVerbose
	}
	if m.Attributes == nil {
		return *oVerbose
	}
	return *oVerbose || len(m.Attributes[attrCloudDebug]) > 0
}
func logInfo(format string, args ...interface{}) {
	logLevel("INFO", nil, format, args...)
}
func logDebug(m *pubsub.Message, format string, args ...interface{}) {
	logLevel("DEBUG", m, format, args...)
}
func logWarn(m *pubsub.Message, format string, args ...interface{}) {
	logLevel("WARN", m, format, args...)
}
func logError(m *pubsub.Message, format string, args ...interface{}) {
	logLevel("ERROR", m, format, args...)
}
func logLevel(level string, m *pubsub.Message, format string, args ...interface{}) {
	log.Printf("[%s] %s: %s\n", tracker(m), level, fmt.Sprintf(format, args...))
}
func tracker(m *pubsub.Message) string {
	if m == nil {
		return "-"
	}
	if m.Attributes == nil {
		return "-"
	}
	t, ok := m.Attributes[attrCloudDebug]
	if !ok {
		return "-"
	}
	if len(t) == 0 {
		return "-"
	}
	if len(m.ID) == 0 {
		return t
	}
	return fmt.Sprintf("%s:%s", t, m.ID)
}
