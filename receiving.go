package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

// loopReceiveParcels receives messages during a limited amount of time (queue.Duration)
// then it enters a sleep to delay receiving the next
func loopReceiveParcels(client *pubsub.Client, queue Queue, service *deliveryServiceImpl) {
	sub := client.Subscription(queue.Subscription)
	for {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			if *verbose {
				log.Printf("will cancel receiving from [%s] after [%v]\n", queue.Subscription, queue.Duration)
			}
			time.Sleep(queue.Duration)
			// abort the receive
			cancel()
			if *verbose {
				log.Printf("cancelled receiving from [%s]\n", queue.Subscription)
			}
		}()
		if *verbose {
			log.Printf("receiving messages from subscription [%s]\n", queue.Subscription)
		}
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			service.receivedParcel(msg)
		})
		if err != nil {
			log.Printf("Receive error from [%s]:%v", queue.Subscription, err)
		}
		if *verbose {
			log.Printf("delay start receiving from [%s] for [%v]\n", queue.Subscription, queue.Duration)
		}
		time.Sleep(queue.Duration)
	}
}
