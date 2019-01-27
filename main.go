package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
)

var (
	oVerbose = flag.Bool("v", false, "verbose logging")
	oConfig  = flag.String("c", "parzello.yaml", "location of configuration")
	version  = "0.5"
)

func main() {
	flag.Parse()
	logInfo("parzello, the pubsub delay service -- version:%s, verbose:%v", version, *oVerbose)
	config, err := loadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// create pubsub client
	if len(config.Project) == 0 {
		log.Fatal("missing project in config")
	}
	ctx := context.Background()
	pubsubClient, err := pubsub.NewClient(ctx, config.Project)
	if err != nil {
		log.Fatalf("failed to create PubSub client: %v", err)
	}
	datastoreClient, err := datastore.NewClient(ctx, config.Project)
	if err != nil {
		log.Fatalf("failed to create Datastore client: %v", err)
	}
	config.checkDataStoreAccessible(datastoreClient)
	config.checkTopicsAndSubscriptions(pubsubClient)

	service := newDelayService(config, pubsubClient, datastoreClient)
	g := new(sync.WaitGroup)
	g.Add(1)
	go func() {
		logInfo("ready to accept messages from subscription [%s]", config.Subscription)
		if err := service.Accept(ctx); err != nil {
			log.Println("accept failed", err)
		}
		g.Done()
	}()
	// schedule listeners
	for _, each := range config.Queues {
		g.Add(1)
		go func(next Queue) {
			loopReceiveParcels(pubsubClient, next, service)
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
