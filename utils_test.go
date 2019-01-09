package main

import (
	"testing"

	"cloud.google.com/go/pubsub"
)

func TestUpdatePublishCount(t *testing.T) {
	m := new(pubsub.Message)
	m.Attributes = map[string]string{}
	updatePublishCount(m)
	if got, want := m.Attributes["parzello.publishCount"], "1"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
	updatePublishCount(m)
	if got, want := m.Attributes["parzello.publishCount"], "2"; got != want {
		t.Errorf("got %v want %v", got, want)
	}
}
