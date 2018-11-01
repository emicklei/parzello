package main

import (
	"log"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
)

const (
	AttrOriginalMessageID = "parzello.originalMessageID"
	AttrDestinationTopic  = "parzello.destinationTopic"
	AttrPublishAfter      = "parzello.publishAfter"
	AttrPublishCount      = "parzello.publishCount"
	AttrEntryTime         = "parzello.entryTime"
)

func timeFromSecondsString(s string) (time.Time, error) {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Now(), err
	}
	return time.Unix(int64(i), 0), nil
}

func secondsToTime(usec uint64) time.Time {
	return time.Unix(int64(usec), 0)
}

func updatePublishCount(m *pubsub.Message) {
	pcs, ok := m.Attributes[AttrPublishCount]
	if ok {
		pc, err := strconv.Atoi(pcs)
		if err != nil {
			log.Println("failed to parse publishCount, set to 1")
			m.Attributes[AttrPublishCount] = "1"
			return
		}
		m.Attributes[AttrPublishCount] = strconv.Itoa(pc + 1)
		return
	}
	// not set
	m.Attributes[AttrPublishCount] = "1"
}

func timeToSecondsString(t time.Time) string {
	return strconv.FormatInt(t.Unix(), 10)
}
