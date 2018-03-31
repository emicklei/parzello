package main

import (
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
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

func (d *deliveryServiceImpl) Deliver(ctx context.Context, req *DeliverRequest) (*DeliverResponse, error) {
	log.Printf("Deliver request for topic [%s] on or after [%v]\n", req.Envelope.DestinationTopic, time.Unix(int64(req.Envelope.PublishAfter.Seconds), 0))

	// select topic to forward message to
	// TODO
	topic := d.client.Topic("parcello_minute")
	p := new(Parcel)
	p.Envelope = req.Envelope
	data, err := proto.Marshal(p)
	if err != nil {
		return &DeliverResponse{ErrorMessage: fmt.Sprintf(errUnknown, err)}, nil
	}
	msg := &pubsub.Message{
		Data: data,
	}
	topic.Publish(ctx, msg)
	return &DeliverResponse{EnqueueTime: &Timestamp{Seconds: uint64(time.Now().Unix())}}, nil
}

// receivedParcel is called from any of the queue subscription pulls.
func (d *deliveryServiceImpl) receivedParcel(msg *pubsub.Message) {
	// payload of message is a Parcel
	p := new(Parcel)
	err := proto.Unmarshal(msg.Data, p)
	if err != nil {
		log.Println("ERROR", err)
		return
	}
	log.Printf("receivedParcel for topic [%s]\n", p.Envelope.DestinationTopic)

	fmt.Printf("ack parcel message [%v]: %q\n", msg.ID, string(p.Envelope.Payload))
	msg.Ack()
	// TEMP
	err = d.publishParcelMessage(p)
	if err != nil {
		log.Println("ERROR", err)
		return
	}
}

// publishParcelMessage publishes the payload of the parcel to the destination topic.
func (d *deliveryServiceImpl) publishParcelMessage(p *Parcel) error {
	topic := d.client.Topic(p.Envelope.DestinationTopic)
	msg := &pubsub.Message{
		Data: p.Envelope.Payload,
	}
	result := topic.Publish(context.Background(), msg)
	log.Printf("%#v", result)
	return nil
}
