package main

import (
	"log"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
)

const (
	attrDestinationTopic = "parzello.destinationTopic"
	attrPublishAfter     = "parzello.publishAfter"
	attrPublishCount     = "parzello.publishCount"
	attrEntryTime        = "parzello.entryTime"
	// if this property is available in the message attributes then all operations for this message are logged in DEBUG level.
	attrCloudDebug = "X-Cloud-Debug"
)

func publishAfter(m *pubsub.Message) time.Time {
	t, _ := timeFromSecondsString(m.Attributes[attrPublishAfter])
	return t
}

func publishCount(m *pubsub.Message) int {
	pcs, ok := m.Attributes[attrPublishCount]
	if ok {
		pc, err := strconv.Atoi(pcs)
		if err != nil {
			return 1
		}
		return pc
	}
	return 1
}

func entryTime(m *pubsub.Message) time.Time {
	t, _ := timeFromSecondsString(m.Attributes[attrEntryTime])
	return t
}

func timeFromSecondsString(s string) (time.Time, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(int64(i), 0), nil
}

func updatePublishCount(m *pubsub.Message) {
	pcs, ok := m.Attributes[attrPublishCount]
	if ok {
		pc, err := strconv.Atoi(pcs)
		if err != nil {
			log.Println("warning: failed to parse publishCount, set to 1")
			m.Attributes[attrPublishCount] = "1"
			return
		}
		m.Attributes[attrPublishCount] = strconv.Itoa(pc + 1)
		return
	}
	// not set
	m.Attributes[attrPublishCount] = "1"
}

func timeToSecondsString(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}

func needsDatastoreMirror(m *pubsub.Message) bool {
	if m.Attributes == nil {
		return false
	}
	return len(m.Attributes[attrDataStoreLookup]) > 0
}
