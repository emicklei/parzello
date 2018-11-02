package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/pubsub"
)

var oDrain = flag.Bool("d", false, "drain destination")
var oCount = flag.Int("c", 1, "how many message")

// GCP_PROJECT=philemonworks go run main.go

func main() {
	flag.Parse()
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, os.Getenv("GCP_PROJECT"))
	if err != nil {
		log.Fatalf("failed to create PubSub client: %v", err)
	}
	defer client.Close()
	topic := client.Topic("parzello_inbound_topic")

	for i := 0; i < *oCount; i++ {
		dur := time.Duration(1+rand.Intn(10)) * time.Minute
		after := time.Now().Add(dur)
		msg := &pubsub.Message{
			Data: []byte("parzello"),
			Attributes: map[string]string{
				"parzello.destinationTopic": "parzello_destination",
				"parzello.publishAfter":     fmt.Sprintf("%d", after.Unix()),
			},
		}
		r := topic.Publish(ctx, msg)
		id, err := r.Get(ctx)
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Published a message to the topic.", id)
	}

	topic.Stop()

	if *oDrain {
		drainDestination(client)
	}
}

func drainDestination(client *pubsub.Client) {
	sub := client.Subscription("parzello_sink")
	log.Println("draining destination", "parzello_sink")
	log.Fatal(sub.Receive(context.Background(), func(ctx context.Context, msg *pubsub.Message) {
		log.Printf(`message-id:%s
-----------------
           parzello.entryTime: %s
                 ps published: %s
                      payload: %s
        parzello.publishAfter: %s
                 actual after: %s
    parzello.destinationTopic: %s
   parzello.originalMessageID: %s
`, msg.ID,
			timeFromSecondsString(msg.Attributes["parzello.entryTime"]).String(),
			msg.PublishTime.String(),
			string(msg.Data),
			timeFromSecondsString(msg.Attributes["parzello.publishAfter"]).String(),
			time.Now().String(),
			msg.Attributes["parzello.destinationTopic"],
			msg.Attributes["parzello.originalMessageID"],
		)
		msg.Ack()
	}))
}

func timeFromSecondsString(s string) time.Time {
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return time.Now()
	}
	return time.Unix(int64(i), 0)
}
