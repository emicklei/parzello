# Example parcello

## prepare infrastructure

    gcloud pubsub topics create parcello_destination
    gcloud pubsub topics create parcello_minute
    gcloud pubsub topics create parcello_five_minutes

and the subscriptions

    gcloud pubsub subscriptions create parcello_destination  --topic parcello_destination --topic-project `gcloud config get-value project`
    gcloud pubsub subscriptions create parcello_minute       --topic parcello_minute      --topic-project `gcloud config get-value project`
    gcloud pubsub subscriptions create parcello_five_minutes --topic parcello_five_minutes --topic-project `gcloud config get-value project`

## start the server

    go run *.go

## running the example

    cd test && go run *.go

## pull from destination

    gcloud pubsub subscriptions pull --auto-ack parcello_destination