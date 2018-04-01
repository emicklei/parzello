package main

import (
	"time"

	"github.com/emicklei/parcello/v1"
	"github.com/golang/protobuf/proto"
)

type DeliverRequestBuilder struct {
	Request *v1.DeliverRequest
	Error   error
}

func NewDeliverRequestBuilder() *DeliverRequestBuilder {
	return &DeliverRequestBuilder{
		Request: &v1.DeliverRequest{
			Envelope: new(v1.Envelope),
		},
	}
}

func (b *DeliverRequestBuilder) PayloadBytes(payload []byte) *DeliverRequestBuilder {
	b.Request.Envelope.Payload = payload
	return b
}

func (b *DeliverRequestBuilder) PayloadMessage(m proto.Message) *DeliverRequestBuilder {
	data, err := proto.Marshal(m)
	if err != nil {
		b.Error = err
	}
	b.Request.Envelope.Payload = data
	return b
}

func (b *DeliverRequestBuilder) Topic(destination string) *DeliverRequestBuilder {
	b.Request.Envelope.DestinationTopic = destination
	return b
}

func (b *DeliverRequestBuilder) Delay(duration time.Duration) *DeliverRequestBuilder {
	b.Request.Envelope.PublishAfter = &v1.Timestamp{Seconds: uint64(time.Now().Add(duration).Unix())}
	return b
}
