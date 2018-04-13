# parcello

Parcello is a delivery service on top of Google Pub Sub to publish messages with a time delay to a topic.
It uses multiple intermediate topic->subscriptions to queue such messages which are pulled (streaming) on computed time intervals.
These intermediate resources (topics, subscription) need to be created before deploying this service with the corresponding configuration.

## use cases
Parcello was created to support a couple of unusual usecases.

### backoff or postpone processing
Consider the situation that the subscriber is unable to process a message because its backend services are not available.
Each such message can not be acknowledge.
Therefore it will be offered on the next pull (or push) causing many operation calls which are responsible for increased network traffic.
One solution is to deliver those messages to `parcello` for the purpose of retrying it later.
Alternatively, the subscriber could have the capability to stop|pause the pulling process.

### late night or bulk processing
Consider the requirement that messages need to be processed in a certain time window only (at night when more resources are available).
One solution is to deliver those messages to `parcello` for the purpose of having them published later.

### out of order
Pub Sub, like most message queueing solutions, does not guarantee that messages are delivered in the same order as they were published. This must be taken into account when designing the message definition and the subscriber software architecture. So rather then publishing events such as "order created", "order detail changed" one should have a single event "order updated" with all information about the order and the timestamp when it was updated.
If this design is not possible then the subscriber must handle the situation that an "order changed" event is received before an "order created" for the same entity.
One solution is to deliver those messages to `parcello` for the purpose of having them published again later.
Ofcourse this is a fragile design and can lead to problems (what if the created event was never published?) so be aware of this when using `parcello` for this case.

## how does it work?

Instead of the publisher publishing to a topic, it sends a DeliverRequest to `parcello` specifying the time to delay the publish. 
The service will publish that request to one of its intermediary topics. 
Each such topic has its own listener (pulling messaging) at a time interval.
In this diagram, 2 such topics exists. One is pulled every 5 minutes. 
Each message is inspected to see whether it is about time to publish it to the actual destination topic.

For example, a message that need to be delivered after 8 minutes will be published to `parcello_5_minutes` once and to `parcello_1_minute ` 3 times.
Using a configuration, you can specify which intermediary topics you have created with durations based on your estimations of the actual delay amounts.

![](./doc/parcello_delay.png)

In the next design, `parcello` is used to retry publishing a message at a later moment.
By passing the PublishCount into the retry DeliverRequest, a subscriber can inspect this value and behave accordingly (i.e abort on MaxRetries).

![](./doc/parcello_delay_retry.png)

## usage

### api

This services uses gRPC for its API access.
The proto definition can be found in the /v1 folder.
In the example folder, you can see an example of a Go client implementation. Using `protoc` tools, you can generate a client for your programming language.

### server config

The `parcello` server configuration list one or more queues (intermediate topic+subscriptions).
This example show 2 such queues. 
The topics and subscriptions must be created upfront in the project `YOUR GCP PROJECT`.
A duration string is a possibly signed sequence of decimal numbers, each with optional fraction and a unit suffix, such as "60s", "1.5h" or "2h45m". Valid time units are "ms", "s", "m", "h".

    {
        "project":"YOUR GCP PROJECT",
        "queues": [
            {
                "topic": "parcello_5_minutes",
                "subscription": "parcello_5_minutes",
                "duration": "5m"
            },
            {
                "topic": "parcello_30_minutes",
                "subscription": "parcello_30_minutes",
                "duration": "30m"
            }
        ]
    }

## build

    docker build -t parcello .

## run

    docker run -it \
        -v ~/.config/gcloud/:/gcloud \
        -e GOOGLE_APPLICATION_CREDENTIALS=/gcloud/application_default_credentials.json \
        parcello