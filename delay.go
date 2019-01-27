package main

import (
	"fmt"
	"sync"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	context "golang.org/x/net/context"
)

type delayService struct {
	pubsubClient    *pubsub.Client
	datastoreClient *datastore.Client
	config          Config
	topicsMutex     *sync.RWMutex
	topics          map[string]*pubsub.Topic
}

func newDelayService(con Config, c *pubsub.Client, d *datastore.Client) *delayService {
	return &delayService{
		pubsubClient:    c,
		datastoreClient: d,
		config:          con,
		topicsMutex:     new(sync.RWMutex),
		topics:          map[string]*pubsub.Topic{},
	}
}

// Accept will start receiving from the public subscription.
func (d *delayService) Accept(ctx context.Context) error {
	sub := d.pubsubClient.Subscription(d.config.Subscription)
	return sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		// validate after
		after, err := timeFromSecondsString(msg.Attributes[attrPublishAfter])
		if err != nil {
			logWarn(msg, "passthrough message because of missing or invalid attribute %s=%s error:%v\n",
				attrPublishAfter, msg.Attributes[attrPublishAfter], err)
			d.publishToDestination(ctx, msg)
			msg.Ack()
			return
		}
		// validate destination
		destination := msg.Attributes[attrDestinationTopic]
		if len(destination) == 0 {
			logError(msg, "unable to handle message because of missing attribute %s=%s",
				attrDestinationTopic, destination)
			msg.Nack()
			return
		}
		if isVerbose(msg) {
			logDebug(msg, "accepted message from subscription [%s] to be delivered to [%s] on or after [%v]",
				d.config.Subscription, destination, after)
		}
		msg.Attributes[attrOriginalMessageID] = msg.ID
		msg.Attributes[attrEntryTime] = timeToSecondsString(msg.PublishTime)
		if err := d.transportMessage(ctx, msg); err != nil {
			logError(msg, "message cannot be transported with error:%v", err)
			msg.Nack()
			return
		}
		msg.Ack()
	})
}

// transportMessage is called from any of the queue subscription pulls or from Deliver.
func (d *delayService) transportMessage(ctx context.Context, m *pubsub.Message) error {
	now := time.Now()
	after, err := timeFromSecondsString(m.Attributes[attrPublishAfter])
	if err != nil {
		return fmt.Errorf("invalid publish after attribute:%v", err)
	}
	// see if it time to publish to the destination
	if after.Before(now) {
		return d.publishToDestination(ctx, m)
	}
	wait := after.Sub(now)
	// pick the queue with the largest duration and within wait
	// at least one exists, has been checked at startup
	nextQueue := d.config.Queues[0]
	for _, each := range d.config.Queues {
		if wait < each.Duration {
			break
		}
		nextQueue = each
	}
	// new message
	msg := &pubsub.Message{
		Data:       m.Data,
		Attributes: m.Attributes,
	}
	if isVerbose(msg) {
		logDebug(msg, "publish message [%s] to [%s]", msg.Attributes[attrOriginalMessageID], nextQueue.Topic)
	}
	d.topicNamed(nextQueue.Topic).Publish(ctx, msg)
	// if the message wants to be mirrored in DataStore then save it
	// unless it was already saved to DataStore upon entry
	if !needsDatastoreMirror(msg) {
		return nil // no need to mirror
	}
	if datastoreKey(msg) != nil {
		return nil // already stored
	}
	k := newDatastoreKey(msg)
	_, err = d.datastoreClient.Put(ctx, k, newPubSubMessageRecord(msg))
	if isVerbose(msg) {
		logDebug(msg, "add mirror copy to datastore with key [%s]", k.Name)
	}
	if err != nil {
		logWarn(msg, "unable to mirror message because [%v]", err)
	}
	return nil
}

// publishToDestination publishes the message to the destination topic.
// if the message was stored in DataStore then remove it.
func (d *delayService) publishToDestination(ctx context.Context, m *pubsub.Message) error {
	// on Accept, the destination has been validated
	destination := m.Attributes[attrDestinationTopic]
	msg := &pubsub.Message{
		Data:       m.Data,
		Attributes: m.Attributes,
	}
	updatePublishCount(msg)
	if isVerbose(msg) {
		msg.ID = msg.Attributes[attrOriginalMessageID]
		logDebug(msg, "publish message to [%s]", destination)
	}
	d.topicNamed(destination).Publish(ctx, msg)
	// if in DataStore then delete it
	if key := datastoreKey(msg); key != nil {
		if err := d.datastoreClient.Delete(ctx, key); err != nil {
			logWarn(msg, "unable to delete datastore entry because [%v]", err)
		} else {
			if isVerbose(msg) {
				logDebug(msg, "deleted message [%s] from datastore", key.Name)
			}
		}
	}
	return nil
}

// topicNamed returns a (cached) pubsub.Topic
func (d *delayService) topicNamed(name string) *pubsub.Topic {
	d.topicsMutex.RLock()
	t, ok := d.topics[name]
	d.topicsMutex.RUnlock()
	if ok {
		return t
	}
	d.topicsMutex.Lock()
	t = d.pubsubClient.Topic(name)
	d.topics[name] = t
	d.topicsMutex.Unlock()
	return t
}
