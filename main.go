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
	oConfig  = flag.String("c", "parzello.yaml", "location of configuration")
	version  = "0.3"
)

func main() {
	flag.Parse()
	log.Printf("parzello, the pubsub delay service -- version:%s, vverbose:%v\n", version, *oVerbose)
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

	service := newDelayService(config, client)
	g := new(sync.WaitGroup)
	g.Add(1)
	go func() {
		log.Printf("ready to accept deliveries on [%s]\n", config.Subscription)
		if err := service.Accept(ctx); err != nil {
			log.Println("accept failed", err)
		}
		g.Done()
	}()
	// schedule listeners
	for _, each := range config.Queues {
		g.Add(1)
		go func(next Queue) {
			loopReceiveParcels(client, next, service)
			g.Done()
		}(each)
	}
	// add selfdiagnose
	g.Add(1)
	go func() {
		addSelfdiagnose()
		g.Done()
	}()
	g.Wait()
}
