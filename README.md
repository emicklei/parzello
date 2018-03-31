# parcello

Parcello is a delivery service on top of Google Pub Sub to publish messages with a time delay and limited retries to a topic.


## prepare infrastructure for example
 
    gcloud pubsub topics create parcello_destination
    gcloud pubsub topics create parcello_undelivered
    gcloud pubsub topics create parcello_minute
    gcloud pubsub topics create parcello_five_minutes

and the subscriptions

    gcloud pubsub subscriptions create parcello_destination  --topic parcello_destination --topic-project philemonworks
    gcloud pubsub subscriptions create parcello_undelivered  --topic parcello_undelivered --topic-project philemonworks
    gcloud pubsub subscriptions create parcello_minute       --topic parcello_minute      --topic-project philemonworks
    gcloud pubsub subscriptions create parcello_five_minutes --topic parcello_five_minutes --topic-project philemonworks

## start the server    

    go run *.go

## running the example

    cd test && go run *.go

## get destination

    gcloud pubsub subscriptions pull --auto-ack parcello_destination
