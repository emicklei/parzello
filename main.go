package main

import (
	"context"
	"log"
	"net"

	"cloud.google.com/go/pubsub"
	"google.golang.org/grpc"
)

func main() {
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	lis, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, "philemonworks")
	if err != nil {
		log.Fatalf("failed to create PubSub client: %v", err)
	}

	for _, each := range config.Queues {
		ok, err := client.Topic(each.Topic).Exists(ctx)
		if err != nil {
			log.Fatalf("failed to check existance: %v", err)
		}
		log.Println("check exists topic", each.Topic)
		if !ok {
			log.Fatalf("topic [%s] does not exist: %v", each.Topic, err)
		}
		log.Println("check exists subscription", each.Subscription)
		ok, err = client.Subscription(each.Subscription).Exists(ctx)
		if err != nil {
			log.Fatalf("failed to check existance: %v", err)
		}
		if !ok {
			log.Fatalf("subscription [%s] does not exist: %v", each.Subscription, err)
		}
	}

	service := &deliveryServiceImpl{client: client, config: config}
	for _, each := range config.Queues {
		go loopReceiveParcels(client, each, service)
	}

	grpcServer := grpc.NewServer()
	RegisterDeliveryServiceServer(grpcServer, service)
	log.Fatal(grpcServer.Serve(lis))
}
