package main

import (
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/emicklei/parcello/v1"
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
	if *oVerbose {
		log.Printf("deliver request for topic [%s] on or after [%v]\n", req.Envelope.DestinationTopic, secondsToTime(req.Envelope.PublishAfter))
	}
	req.Envelope.DeliveredAt = uint64(time.Now().Unix())
	req.Envelope.ID = newUUID()
	if err := validateEnvelop(req.Envelope); err != nil {
		return &v1.DeliverResponse{ErrorMessage: err.Error()}, nil
	}
	if err := d.transportParcel(ctx, req.Envelope); err != nil {
		return &v1.DeliverResponse{ErrorMessage: fmt.Sprintf(errUnknown, err)}, nil
	}

	return new(v1.DeliverResponse), nil
}

// transportParcel is called from any of the queue subscription pulls or from Deliver.
func (d *deliveryServiceImpl) transportParcel(ctx context.Context, e *v1.Envelope) error {
	now := time.Now()
	after := secondsToTime(e.PublishAfter)

	// see if destination is arrived
	if after.Before(now) {
		return d.publishMessageOfParcel(e)
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
	data, err := proto.Marshal(e)
	if err != nil {
		return err
	}
	if *oVerbose {
		log.Printf("publish parcel [%s] to [%s]", e.ID, nextQueue.Topic)
	}
	d.client.Topic(nextQueue.Topic).Publish(ctx, &pubsub.Message{Data: data})
	return nil
}

// publishMessageOfParcel publishes the payload of the parcel to the destination topic.
func (d *deliveryServiceImpl) publishMessageOfParcel(e *v1.Envelope) error {
	if *oVerbose {
		log.Printf("publish message from parcel [%s] to [%s]", e.ID, e.DestinationTopic)
	}
	topic := d.client.Topic(e.DestinationTopic)
	msg := &pubsub.Message{
		Data:       e.Payload,
		Attributes: e.Attributes,
	}
	setMessageAttributes(e, msg)
	topic.Publish(context.Background(), msg)
	return nil
}
