package main

import (
	"context"
	"log"
	"os"
	"sort"
	"time"

	"cloud.google.com/go/pubsub"
	yaml "gopkg.in/yaml.v3"
)

// Queue is a topic,subscription pair for storing Parcels for some duration.
type Queue struct {
	Topic             string        `yaml:"topic"`
	Subscription      string        `yaml:"subscription"`
	FormattedDuration string        `yaml:"duration"`
	Duration          time.Duration `yaml:"-"`
}

func (q Queue) postLoaded() (Queue, error) {
	dur, err := time.ParseDuration(q.FormattedDuration)
	return Queue{Topic: q.Topic, Subscription: q.Subscription, Duration: dur}, err
}

// Config is read from a configuration JSON file.
type Config struct {
	Project      string `yaml:"project-id"`
	Subscription string `yaml:"subscription"`
	// Queues is sorted by Duration, shortest first
	Queues []Queue `yaml:"queues"`
}

func (c Config) checkTopicsAndSubscriptions(client *pubsub.Client) {
	if len(c.Queues) == 0 {
		log.Fatalln("at least one queue (topic,subscription) must be configured")
	}
	ctx := context.Background()
	// check inbound subscription
	ok, err := client.Subscription(c.Subscription).Exists(ctx)
	if err != nil {
		log.Fatalf("failed to check existence: %v", err)
	}
	if !ok {
		log.Fatalf("subscription [%s] does not exist", c.Subscription)
	}
	// check all queues
	for _, each := range c.Queues {
		ok, err := client.Topic(each.Topic).Exists(ctx)
		if err != nil {
			log.Fatalf("failed to check existence: %v", err)
		}
		log.Printf("check exists topic [%s]\n", each.Topic)
		if !ok {
			log.Fatalf("topic [%s] does not exist: %v", each.Topic, err)
		}
		log.Printf("check exists subscription [%s]\n", each.Subscription)
		ok, err = client.Subscription(each.Subscription).Exists(ctx)
		if err != nil {
			log.Fatalf("failed to check existence: %v", err)
		}
		if !ok {
			log.Fatalf("subscription [%s] does not exist", each.Subscription)
		}
	}
}

func loadConfig() (config Config, err error) {
	data, err := os.ReadFile(*oConfig)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &config)
	for i, each := range config.Queues {
		other, err := each.postLoaded()
		if err != nil {
			return config, err
		}
		config.Queues[i] = other
	}
	// sort by duration ascending; shortest first
	sort.Slice(config.Queues, func(i, j int) bool { return config.Queues[i].Duration < config.Queues[j].Duration })
	log.Println("loaded configuration from", *oConfig)
	return config, err
}
