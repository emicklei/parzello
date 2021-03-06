# Example parzello

## prepare infrastructure

    gcloud pubsub topics create parzello_inbound_topic
    gcloud pubsub topics create parzello_destination
    gcloud pubsub topics create parzello_minute
    gcloud pubsub topics create parzello_five_minutes

and the subscriptions

    gcloud pubsub subscriptions create parzello_destination  --topic parzello_destination --topic-project `gcloud config get-value project`
    gcloud pubsub subscriptions create parzello_minute       --topic parzello_minute      --topic-project `gcloud config get-value project`
    gcloud pubsub subscriptions create parzello_five_minutes --topic parzello_five_minutes --topic-project `gcloud config get-value project`
    gcloud pubsub subscriptions create parzello_inbound --topic parzello_inbound_topic --topic-project `gcloud config get-value project`

## start the server
In the project root

    make run

## running the example
In the example folder

    GCP=`gcloud config get-value project` go run *.go

## pull from destination

    gcloud pubsub subscriptions pull --auto-ack parzello_destination
