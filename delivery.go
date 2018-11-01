package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/pubsub"
	context "golang.org/x/net/context"
)

var (
	errUnknown = "1000: unknown error occurred [%v]"
	errPublish = "1001: failed to publish to [%v]"
)

type deliveryService struct {
	client      *pubsub.Client
	config      Config
	topicsMutex *sync.RWMutex
	topics      map[string]*pubsub.Topic
}

func newDeliveryService(con Config, c *pubsub.Client) *deliveryService {
	return &deliveryService{
		client:      c,
		config:      con,
		topicsMutex: new(sync.RWMutex),
		topics:      map[string]*pubsub.Topic{},
	}
}

func (d *deliveryService) Accept(ctx context.Context) error {
	sub := d.client.Subscription(d.config.Subscription)
	return sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		if *oVerbose {
			// TODO handle err
			t, _ := timeFromSecondsString(msg.Attributes[AttrPublishAfter])
			log.Printf("accept message from topic [%s] to be delivered on or after [%v]\n", msg.Attributes[AttrDestinationTopic], t)
		}
		msg.Attributes[AttrOriginalMessageID] = msg.ID
		msg.Attributes[AttrEntryTime] = timeToSecondsString(msg.PublishTime)
		if err := d.transportMessage(ctx, msg); err != nil {
			log.Printf("unable to transport message:%v\n", err)
			return
		}
		if *oVerbose {
			log.Println("acknowledge accepted message", msg.ID)
		}
		msg.Ack()
	})
}

// transportMessage is called from any of the queue subscription pulls or from Deliver.
func (d *deliveryService) transportMessage(ctx context.Context, m *pubsub.Message) error {
	now := time.Now()
	after, err := timeFromSecondsString(m.Attributes[AttrPublishAfter])
	if err != nil {
		return fmt.Errorf("invalid publish after attribute:%v", err)
	}

	// see if destination is arrived
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
	// TODO set PublishCount
	// TODO set Time info
	msg := &pubsub.Message{
		Data:       m.Data,
		Attributes: m.Attributes,
	}
	if *oVerbose {
		log.Printf("publish message [%s] to [%s]", m.ID, nextQueue.Topic)
	}
	d.topicNamed(nextQueue.Topic).Publish(ctx, msg)
	return nil
}

// publishToDestination publishes the message to the destination topic.
func (d *deliveryService) publishToDestination(ctx context.Context, m *pubsub.Message) error {
	if *oVerbose {
		log.Printf("publish message [%s] to [%s]", m.ID, m.Attributes[AttrDestinationTopic])
	}
	// new message
	// TODO set PublishCount
	// TODO set Time info
	msg := &pubsub.Message{
		Data:       m.Data,
		Attributes: m.Attributes,
	}
	//setMessageAttributes(e, msg)
	updatePublishCount(m)
	d.topicNamed(m.Attributes[AttrDestinationTopic]).Publish(ctx, msg)
	return nil
}

func (d *deliveryService) topicNamed(name string) *pubsub.Topic {
	d.topicsMutex.RLock()
	t, ok := d.topics[name]
	d.topicsMutex.RUnlock()
	if ok {
		return t
	}
	// create
	d.topicsMutex.Lock()
	t = d.client.Topic(name)
	d.topics[name] = t
	d.topicsMutex.Unlock()
	return t
}
