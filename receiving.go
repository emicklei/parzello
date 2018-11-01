package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

// loopReceiveParcels receives messages during a limited amount of time (queue.Duration)
// then it enters a sleep to delay receiving the next
func loopReceiveParcels(client *pubsub.Client, q Queue, service *deliveryServiceImpl) {
	sub := client.Subscription(q.Subscription)
	for {
		ctx, cancel := context.WithCancel(context.Background())
		if *oVerbose {
			log.Printf("receiving messages from subscription [%s]\n", q.Subscription)
		}
		var delayBeforeReceive time.Duration
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			// if the message was fetched too soon then cancel this receiving and enter a sleep
			now := time.Now()
			after, err := timeFromSecondsString(msg.Attributes[AttrPublishAfter])
			if err != nil {
				log.Printf("invalid publish after attribute:%v\n", err)
				return
			}
			if after.After(now) {
				// compute remaining time for this message to stay in queues
				delayBeforeReceive = after.Sub(now)
				if *oVerbose {
					log.Printf("nack message [%s] and cancel receiving to enter delay, [%v] in subscription [%s], remaining [%v]\n", msg.ID, now.Sub(msg.PublishTime), q.Subscription, delayBeforeReceive)
				}
				msg.Nack()
				cancel()
				return
			}
			if err := service.transportMessage(ctx, msg); err == nil {
				msg.Ack()
			} else {
				log.Println("transportMessage error", err)
				msg.Nack()
			}
		})
		if err != nil {
			log.Printf("receive error from [%s]:%v", q.Subscription, err)
		}
		if delayBeforeReceive > q.Duration {
			delayBeforeReceive = q.Duration // do not wait longer than specified by the queue
		}
		if *oVerbose {
			log.Printf("delay start receiving from [%s] for [%v]\n", q.Subscription, delayBeforeReceive)
		}
		time.Sleep(delayBeforeReceive)
	}
}
