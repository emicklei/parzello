package main

import (
	"context"
	"time"

	"cloud.google.com/go/pubsub"
)

// loopReceiveParcels receives messages during a limited amount of time (queue.Duration)
// then it enters a sleep to delay receiving the next
func loopReceiveParcels(client *pubsub.Client, q Queue, service *delayService) {
	sub := client.Subscription(q.Subscription)
	if *oVerbose {
		logInfo("ready to receive messages from subscription [%s]", q.Subscription)
	}
	for {
		ctx, cancel := context.WithCancel(context.Background())
		delayBeforeReceive := q.Duration
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			now := time.Now()
			after, err := timeFromSecondsString(msg.Attributes[attrPublishAfter])
			if err != nil {
				logError(msg, "message has invalid [%s] attribute:%v\n", attrPublishAfter, err)
				return
			}
			// if the message was fetched too soon then cancel this receive and enter a sleep
			if after.After(now) {
				// compute remaining time for this message to stay in the subscription
				delayBeforeReceive = after.Sub(now)
				if isVerbose(msg) {
					logDebug(msg, "message is rejected and receive is cancelled to enter delay [%v] in subscription [%s]",
						delayBeforeReceive, q.Subscription)
				}
				msg.Nack()
				cancel()
				return
			}
			if err := service.transportMessage(ctx, msg); err == nil {
				msg.Ack()
			} else {
				logError(msg, "message cannot be transported with error:%v", err)
				msg.Nack()
			}
		})
		if err != nil {
			logError(nil, "receive from subscription [%v] failed with error:%v", q.Subscription, err)
		}
		if delayBeforeReceive > q.Duration {
			delayBeforeReceive = q.Duration // do not wait longer than specified by the queue
		}
		if *oVerbose {
			logDebug(nil, "delay start receiving from subscription [%s] for [%v]\n", q.Subscription, delayBeforeReceive)
		}
		time.Sleep(delayBeforeReceive)
	}
}
