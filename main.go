package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"cloud.google.com/go/pubsub"
)

var (
	oVerbose = flag.Bool("v", false, "verbose logging")
	oConfig  = flag.String("c", "parzello.config", "location of configuration")
	version  = "0.2"
)

func main() {
	flag.Parse()
	log.Println("parzello, the pubsub delivery service -- version", version)
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// create pubsub client
	if len(config.Project) == 0 {
		log.Fatal("missing project in config")
	}
	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, config.Project)
	if err != nil {
		log.Fatalf("failed to create PubSub client: %v", err)
	}

	config.checkTopicsAndSubscriptions(client)

	// schedule parcel listeners
	service := &deliveryServiceImpl{client: client, config: config}
	g := new(sync.WaitGroup)
	g.Add(1)
	go func() {
		if err := service.Accept(ctx); err != nil {
			log.Println("Accept failed", err)
		}
		g.Done()
	}()
	log.Println("ready to accept deliveries....")
	for _, each := range config.Queues {
		g.Add(1)
		go func(next Queue) {
			loopReceiveParcels(client, next, service)
			g.Done()
		}(each)
	}
	log.Println("ready to handle deliveries....")
	g.Wait()
}
