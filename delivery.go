package main

import (
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/emicklei/parcello/v1"
	"github.com/emicklei/parcello/v1/queue"
	"github.com/golang/protobuf/proto"
	context "golang.org/x/net/context"
)

var (
	errUnknown = "1000: unknown error occurred [%v]"
	errPublish = "1001: failed to publish to [%v]"
)

type deliveryServiceImpl struct {
	client *pubsub.Client
	config Config
}

func (d *deliveryServiceImpl) Deliver(ctx context.Context, req *v1.DeliverRequest) (*v1.DeliverResponse, error) {
	if *verbose {
		log.Printf("deliver request for topic [%s] on or after [%v]\n", req.Envelope.DestinationTopic, timestampToTime(req.Envelope.PublishAfter))
	}

	p := new(queue.Parcel)
	p.Envelope = req.Envelope
	if err := d.transportParcel(ctx, p); err != nil {
		return &v1.DeliverResponse{ErrorMessage: fmt.Sprintf(errUnknown, err)}, nil
	}

	return new(v1.DeliverResponse), nil
}

// transportParcel is called from any of the queue subscription pulls or from Deliver.
func (d *deliveryServiceImpl) transportParcel(ctx context.Context, p *queue.Parcel) error {
	now := time.Now()
	after := timestampToTime(p.Envelope.PublishAfter)

	// see if destination is arrived
	if after.Before(now) {
		return d.publishMessageOfParcel(p)
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

	// publish parcel to intermediate queue
	topic := d.client.Topic(nextQueue.Topic)
	data, err := proto.Marshal(p)
	if err != nil {
		return err
	}
	if *verbose {
		log.Printf("publish parcel to [%s]", nextQueue.Topic)
	}
	topic.Publish(ctx, &pubsub.Message{Data: data})
	return nil
}

// publishMessageOfParcel publishes the payload of the parcel to the destination topic.
func (d *deliveryServiceImpl) publishMessageOfParcel(p *queue.Parcel) error {
	if *verbose {
		log.Printf("publish message to [%s]", p.Envelope.DestinationTopic)
	}
	topic := d.client.Topic(p.Envelope.DestinationTopic)
	msg := &pubsub.Message{
		Data:       p.Envelope.Payload,
		Attributes: p.Envelope.Attributes,
	}
	topic.Publish(context.Background(), msg)
	return nil
}

func timestampToTime(t *v1.Timestamp) time.Time {
	return time.Unix(int64(t.Seconds), 0)
}
