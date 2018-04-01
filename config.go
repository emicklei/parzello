package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"sort"
	"time"

	"cloud.google.com/go/pubsub"
)

// Queue is a topic,subscription pair for storing Parcels for some duration.
type Queue struct {
	Topic             string `json:"topic"`
	Subscription      string `json:"subscription"`
	FormattedDuration string `json:"duration"`
	Duration          time.Duration
}

func (q Queue) postLoaded() (Queue, error) {
	dur, err := time.ParseDuration(q.FormattedDuration)
	return Queue{Topic: q.Topic, Subscription: q.Subscription, Duration: dur}, err
}

// Config is read from a configuration JSON file.
type Config struct {
	Project string `json:"project"`
	// Queues is sorted by Duration, shortest first
	Queues []Queue `json:"queues"`
}

func (c Config) checkTopicsAndSubscriptions(client *pubsub.Client) {
	ctx := context.Background()
	for _, each := range c.Queues {
		ok, err := client.Topic(each.Topic).Exists(ctx)
		if err != nil {
			log.Fatalf("failed to check existance: %v", err)
		}
		log.Printf("check exists topic [%s]\n", each.Topic)
		if !ok {
			log.Fatalf("topic [%s] does not exist: %v", each.Topic, err)
		}
		log.Printf("check exists subscription [%s]\n", each.Subscription)
		ok, err = client.Subscription(each.Subscription).Exists(ctx)
		if err != nil {
			log.Fatalf("failed to check existance: %v", err)
		}
		if !ok {
			log.Fatalf("subscription [%s] does not exist: %v", each.Subscription, err)
		}
	}
}

func loadConfig() (config Config, err error) {
	data, err := ioutil.ReadFile("parcello.config")
	if err != nil {
		return
	}
	err = json.Unmarshal(data, &config)
	for i, each := range config.Queues {
		other, err := each.postLoaded()
		if err != nil {
			return config, err
		}
		config.Queues[i] = other
	}
	// sort by duration ascending; shortest first
	sort.Slice(config.Queues, func(i, j int) bool { return config.Queues[i].Duration < config.Queues[j].Duration })
	return config, err
}
