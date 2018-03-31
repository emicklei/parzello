package main

import (
	"context"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
)

func loopReceiveParcels(client *pubsub.Client, queue Queue, service *deliveryServiceImpl) {
	sub := client.Subscription(queue.Subscription)
	for {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			log.Printf("cancel receiving from [%s] after [%v]\n", queue.Subscription, queue.Duration)
			time.Sleep(queue.Duration)
			// interrupt the receive
			cancel()
			log.Println("cancelled receiving from", queue.Subscription)
		}()
		log.Println("receiving messages from subscription", queue.Subscription)
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			service.receivedParcel(msg)
		})
		if err != nil {
			log.Printf("Receive error from [%s]:%v", queue.Subscription, err)
		}
		log.Printf("delay start receiving from [%s] for [%v]\n", queue.Subscription, queue.Duration)
		time.Sleep(queue.Duration)
	}
}
