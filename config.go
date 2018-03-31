package main

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

type Config struct {
	Queues []Queue `json:"queues"`
}

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
	// sort by duration TODO
	return config, err
}
