package main

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/emicklei/parcello/v1"
)

func validateEnvelop(e *v1.Envelope) error {
	if len(e.DestinationTopic) == 0 {
		return errors.New("destinationTopic cannot be empty")
	}
	if e.PublishAfter == 0 {
		return errors.New("publishAfter cannot be zero")
	}
	return nil
}

func secondsToTime(usec uint64) time.Time {
	return time.Unix(int64(usec), 0)
}

func setMessageAttributes(e *v1.Envelope, m *pubsub.Message) {
	if m.Attributes == nil {
		m.Attributes = map[string]string{}
	}
	m.Attributes["parcello.ID"] = e.ID
	m.Attributes["parcello.destinationTopic"] = e.DestinationTopic
	m.Attributes["parcello.publishAfter"] = secondsToTime(e.PublishAfter).UTC().String()
	m.Attributes["parcello.deliveredAt"] = secondsToTime(e.DeliveredAt).UTC().String()
}

func newUUID() string {
	u := make([]byte, 16)
	_, _ = rand.Read(u)
	return hex.EncodeToString(u)
}
