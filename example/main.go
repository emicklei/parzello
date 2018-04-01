package main

import (
	"context"
	"log"
	"time"

	parcello "github.com/emicklei/parcello/v1"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Dial failed:", err)
	}
	defer conn.Close()
	client := parcello.NewDeliveryServiceClient(conn)

	after := time.Now().Add(1 * time.Minute)
	in := new(parcello.DeliverRequest)
	in.Envelope = &parcello.Envelope{
		Payload:          []byte("parcello " + time.Now().String()),
		DestinationTopic: "parcello_destination",
		UndeliveredTopic: "parcello_undelivered",
		PublishAfter:     &parcello.Timestamp{Seconds: uint64(after.Unix())},
	}
	out, err := client.Deliver(context.Background(), in)
	if err != nil {
		log.Fatal("Deliver failed:", err)
	}
	log.Printf("%#v", out)
}
