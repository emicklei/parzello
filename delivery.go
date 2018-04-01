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
		log.Printf("Deliver request for topic [%s] on or after [%v]\n", req.Envelope.DestinationTopic, time.Unix(int64(req.Envelope.PublishAfter.Seconds), 0))
	}

	// select topic to forward message to
	// TODO
	topic := d.client.Topic("parcello_minute")
	p := new(queue.Parcel)
	p.Envelope = req.Envelope
	data, err := proto.Marshal(p)
	if err != nil {
		return &v1.DeliverResponse{ErrorMessage: fmt.Sprintf(errUnknown, err)}, nil
	}
	msg := &pubsub.Message{
		Data: data,
	}
	topic.Publish(ctx, msg)
	return &v1.DeliverResponse{EnqueueTime: &v1.Timestamp{Seconds: uint64(time.Now().Unix())}}, nil
}

// receivedParcel is called from any of the queue subscription pulls.
func (d *deliveryServiceImpl) receivedParcel(msg *pubsub.Message) {
	// payload of message is a Parcel
	p := new(queue.Parcel)
	err := proto.Unmarshal(msg.Data, p)
	if err != nil {
		log.Println("ERROR", err)
		return
	}
	if *verbose {
		log.Printf("receivedParcel for topic [%s]\n", p.Envelope.DestinationTopic)
	}

	// TEMP
	err = d.publishParcelMessage(p)
	if err != nil {
		log.Println("ERROR", err)
		return
	}

	if *verbose {
		log.Printf("ack parcel [%v]: %q\n", msg.ID, string(p.Envelope.Payload))
	}
	msg.Ack()
}

// publishParcelMessage publishes the payload of the parcel to the destination topic.
func (d *deliveryServiceImpl) publishParcelMessage(p *queue.Parcel) error {
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
