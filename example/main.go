package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
	parzello "github.com/emicklei/parzello/v1"
	"google.golang.org/grpc"
)

// GCP_PROJECT=philemonworks go run main.go

func main() {
	conn, err := grpc.Dial("localhost:9090", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Dial failed:", err)
	}
	defer conn.Close()

	go drainDestination()

	client := parzello.NewDeliveryServiceClient(conn)

	for i := 0; i < 100; i++ {
		d := time.Duration(rand.Intn(12)) * time.Minute
		//d, _ := time.ParseDuration("1m30s")
		after := time.Now().Add(d)
		in := new(parzello.DeliverRequest)
		in.Envelope = &parzello.Envelope{
			Payload:          []byte(strconv.Itoa(i)),
			DestinationTopic: "parzello_destination",
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
	client, err := pubsub.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		log.Fatalf("failed to create PubSub client: %v", err)
	}
	sub := client.Subscription("parzello_destination")
	err = sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
		log.Printf(`id:%s
-----------------			
         parzello.deliveredAt:%s
                 ps published:%s
                      payload:%s
        parzello.publishAfter:%s
                 actual after:%s
    parzello.destinationTopic:%s
                  parzello.ID:%s
`, msg.ID,
			msg.Attributes["parzello.deliveredAt"],
			msg.PublishTime.UTC().String(),
			string(msg.Data),
			msg.Attributes["parzello.publishAfter"],
			time.Now().UTC().String(),
			msg.Attributes["parzello.destinationTopic"],
			msg.Attributes["parzello.ID"],
		)
		msg.Ack()
	})
	if err != nil {
		log.Printf("Receiving stopped with error %v", err)
	}
}
