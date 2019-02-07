package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

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

	// start inbound receive process
	service := newDelayService(config, pubsubClient, datastoreClient)
	go func() {
		logInfo("ready to accept messages from subscription [%s]", config.Subscription)
		if err := service.Accept(ctx); err != nil {
			log.Println("accept failed", err)
		}
	}()
	// start internal receiver processes
	for _, each := range config.Queues {
		go func(next Queue) {
			loopReceiveParcels(pubsubClient, next, service)
		}(each)
	}
	addSelfdiagnose()
	addAPI(datastoreClient, config)
	startHTTP()
}

func addAPI(d *datastore.Client, c Config) {
	api := NewAPI(d, c)
	http.HandleFunc("/v1/count", api.counts)
}

func startHTTP() {
	// start a HTTP server
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/internal/selfdiagnose.html", http.StatusTemporaryRedirect)
	})

	webServer := &http.Server{Addr: ":" + port}
	logInfo("HTTP api is listening on %s", port)
	logInfo("[DEV] open http://localhost:%s/internal/selfdiagnose.html", port)
	webServer.ListenAndServe()
}
