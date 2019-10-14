package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"time"

	"cloud.google.com/go/datastore"
	"cloud.google.com/go/pubsub"
	selfdiagnose "github.com/emicklei/go-selfdiagnose"
	"github.com/go-yaml/yaml"
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
	Queues    []Queue     `yaml:"queues"`
	BasicAuth Credentials `yaml:"basic-auth"`
}

type Credentials struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
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
	// register task for future detection
	selfdiagnose.Register(&CheckSubscriptionTask{Client: client, Subscription: c.Subscription})

	// check all queues
	for _, each := range c.Queues {
		ok, err := client.Topic(each.Topic).Exists(ctx)
		if err != nil {
			log.Fatalf("failed to check existence: %v", err)
		}
		logInfo("check exists topic [%s]", each.Topic)
		if !ok {
			log.Fatalf("topic [%s] does not exist: %v", each.Topic, err)
		}
		logInfo("check exists subscription [%s]", each.Subscription)
		ok, err = client.Subscription(each.Subscription).Exists(ctx)
		if err != nil {
			log.Fatalf("failed to check existence: %v", err)
		}
		if !ok {
			log.Fatalf("subscription [%s] does not exist", each.Subscription)
		}
		// register task for future detection
		selfdiagnose.Register(&CheckTopicTask{Client: client, Topic: each.Topic})
		selfdiagnose.Register(&CheckSubscriptionTask{Client: client, Subscription: each.Subscription})
		selfdiagnose.Register(selfdiagnose.ReportMessage{Message: fmt.Sprintf("subscription [%s] is pulled every [%v]", each.Subscription, each.Duration)})
	}
}

func (c Config) checkDataStoreAccessible(datastoreClient *datastore.Client) {
	countQuery := datastore.NewQuery("__Stat_Ns_Kind__")
	countQuery.Namespace("parzello")
	countQuery.Filter("kind_name =", kindDatastore)
	n, err := datastoreClient.Count(context.Background(), countQuery)
	if err != nil {
		log.Fatalf("count on parzello.%s failed [%v]", kindDatastore, err)
	}
	logInfo("[%d] parzello.%s pending", n, kindDatastore)
}

func loadConfig() (config Config, err error) {
	data, err := ioutil.ReadFile(*oConfig)
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
	logInfo("loaded configuration from %v", *oConfig)
	return config, err
}
