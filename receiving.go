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
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second) // TODO config
		delayBeforeReceive := q.Duration
		err := sub.Receive(ctx, func(ctx context.Context, msg *pubsub.Message) {
			now := time.Now()
			after, err := timeFromSecondsString(msg.Attributes[attrPublishAfter])
			if err != nil {
				logWarn(msg, "cannot delay message because it has invalid [%s] attribute:%v", attrPublishAfter, err)
				// allow message to be transported to destination
				after = now.Add(1 * time.Second)
			}
			// if the message was fetched too soon then cancel this receive and enter a sleep
			if after.After(now) {
				// compute remaining time for this message to stay in the subscription
				delayBeforeReceive = after.Sub(now)
				if isVerbose(msg) {
					logDebug(msg, "receive is cancelled to enter delay [%v] in subscription [%s]",
						delayBeforeReceive, q.Subscription)
				}
				msg.Nack()
				cancel() // abort receive, enter sleep
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
			logDebug(nil, "delay start receiving from subscription [%s] for [%v]", q.Subscription, delayBeforeReceive)
		}
		time.Sleep(delayBeforeReceive)
	}
}
