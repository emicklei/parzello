package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/emicklei/parcello/v1/queue"
	"github.com/golang/protobuf/proto"
)

// loopReceiveParcels receives messages during a limited amount of time (queue.Duration)
// then it enters a sleep to delay receiving the next
func loopReceiveParcels(client *pubsub.Client, q Queue, service *deliveryServiceImpl) {
	sub := client.Subscription(q.Subscription)
	for {
		ctx, cancel := context.WithCancel(context.Background())
		if *verbose {
			log.Printf("receiving messages from subscription [%s]\n", q.Subscription)
		}
		var delayBeforeReceive time.Duration
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			// payload of message is a Parcel
			p := new(queue.Parcel)
			err := proto.Unmarshal(msg.Data, p)
			if err != nil {
				log.Println("Unmarshal parcel error", err)
				msg.Nack()
				return
			}
			// if the message was fetched too soon then cancel this receiving and enter a sleep
			now := time.Now()
			if msg.PublishTime.Add(q.Duration).After(now) {
				// compute remaining time for this message to stay in the queue
				delayBeforeReceive = msg.PublishTime.Add(q.Duration).Sub(now)
				if *verbose {
					log.Printf("nack parcel and cancel receiving to enter delay, [%v] in queue, remaining [%v]\n", now.Sub(msg.PublishTime), delayBeforeReceive)
				}
				msg.Nack()
				cancel()
				return
			}
			if err := service.transportParcel(ctx, p); err == nil {
				msg.Ack()
			} else {
				log.Println("transportParcel error", err)
				msg.Nack()
			}
		})
		if err != nil {
			log.Printf("Receive error from [%s]:%v", q.Subscription, err)
		}
		if delayBeforeReceive > q.Duration {
			delayBeforeReceive = q.Duration // do not wait langer than specified by the queue
		}
		if *verbose {
			log.Printf("delay start receiving from [%s] for [%v]\n", q.Subscription, delayBeforeReceive)
		}
		time.Sleep(delayBeforeReceive)
	}
}
