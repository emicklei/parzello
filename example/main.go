package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
	parcello "github.com/emicklei/parcello/v1"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Dial failed:", err)
	}
	defer conn.Close()

	go drainDestination()

	client := parcello.NewDeliveryServiceClient(conn)

	for i := 0; i < 1; i++ {
		//d := time.Duration(rand.Intn(10)) * time.Minute
		d, _ := time.ParseDuration("1m30s")
		after := time.Now().Add(d)
		in := new(parcello.DeliverRequest)
		in.Envelope = &parcello.Envelope{
			Payload:          []byte(strconv.Itoa(i)),
			DestinationTopic: "parcello_destination",
			PublishAfter:     uint64(after.Unix()),
		}
		out, err := client.Deliver(context.Background(), in)
		if err != nil {
			log.Fatal("Deliver failed:", err)
		}
		log.Printf("%#v", out)
	}
	fmt.Print("enter to exit ... ")
	bufio.NewReader(os.Stdin).ReadString('\n')
}

func drainDestination() {
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "philemonworks")
	if err != nil {
		log.Fatalf("failed to create PubSub client: %v", err)
	}
	sub := client.Subscription("parcello_destination")
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Printf(`id:%s
                 ps published:%s
                      payload:%s
                  parcello.ID:%s
        parcello.publishAfter:%s
                 actual after:%s		
    parcello.destinationTopic:%s
         parcello.deliveredAt:%s
`, msg.ID,
			msg.PublishTime.UTC().String(),
			string(msg.Data),
			msg.Attributes["parcello.ID"],
			msg.Attributes["parcello.publishAfter"],
			time.Now().UTC().String(),
			msg.Attributes["parcello.destinationTopic"],
			msg.Attributes["parcello.deliveredAt"],
		)
		msg.Ack()
	})
	if err != nil {
		log.Printf("Receiving stopped with error %v", err)
	}
}
