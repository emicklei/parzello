package main

import (
	"encoding/json"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	"github.com/google/uuid"
)

const (
	attrDataStoreMirror = "parzello.datastoreMirror" // true or false
	attrDataStoreInfo   = "parzello.datastoreInfo"   // set by Publisher to allow query
	attrDataStoreKey    = "parzello.datastoreKey"    // set by Parzello for delete

	kindDatastore = "PubSubRecord"
)

type PubSubRecord struct {
	Info             string
	PublishAfter     time.Time
	DestinationTopic string
	PublishCount     int
	EntryTime        time.Time
	Payload          []byte
	Properties       string // JSON representation
}

func newPubSubMessageRecord(m *pubsub.Message) *PubSubRecord {
	data, _ := json.Marshal(m.Attributes)
	return &PubSubRecord{
		Info:             m.Attributes[attrDataStoreInfo],
		PublishAfter:     publishAfter(m),
		DestinationTopic: m.Attributes[attrDestinationTopic],
		PublishCount:     publishCount(m),
		EntryTime:        entryTime(m),
		Payload:          m.Data,
		Properties:       string(data),
	}
}

// Return key or empty string
func datastoreKey(m *pubsub.Message) *datastore.Key {
	if m.Attributes == nil {
		return nil
	}
	id := m.Attributes[attrDataStoreKey]
	if len(id) == 0 {
		return nil
	}
	k := datastore.NameKey(kindDatastore, id, nil)
	k.Namespace = "parzello"
	return k
}

// newDatastoreKey returns a new datastore.Key and stores the key in the message attributes
func newDatastoreKey(m *pubsub.Message) *datastore.Key {
	if m.Attributes == nil {
		m.Attributes = map[string]string{}
	}
	id := uuid.New().String()
	m.Attributes[attrDataStoreKey] = id
	k := datastore.NameKey(kindDatastore, id, nil)
	k.Namespace = "parzello"
	return k
}
