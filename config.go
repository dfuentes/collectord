package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	Sinks    []map[string]string `json:"sinks"`
	Sources  []map[string]string `json:"sources"`
	Channels []map[string]string `json:"channels"`
	Location string
}

type ComponentSettings map[string]string

var config Config

var sinkLookup map[string]Sink
var channelLookup map[string]Channel
var sourceLookup map[string]Source

func init() {
	confUsage := fmt.Sprintf("Set the config file.  This can also be set by the environment variable %s", CONFIG_ENV)
	flag.StringVar(&config.Location, "conf", "", confUsage)
}

func SetupConfig() {
	if config.Location == "" {
		// No config specified on command line, try environment variable
		config.Location = os.Getenv(CONFIG_ENV)
	}
	if config.Location == "" {
		// env variable also not set
		log.Fatal("No config location specified")
	}

	if _, err := os.Stat(config.Location); os.IsNotExist(err) {
		log.Fatalf("Config file does not exist: %s", config.Location)
	}

	loadConfig()
}

func loadConfig() {
	sinkLookup = make(map[string]Sink)
	channelLookup = make(map[string]Channel)
	sourceLookup = make(map[string]Source)

	rawConfig, err := ioutil.ReadFile(config.Location)
	if err != nil {
		log.Fatalf("Could not open config file for reading: %s", err)
	}
	if err = json.Unmarshal(rawConfig, &config); err != nil {
		log.Fatalf("Error reading config json: %s", err)
	}

	// init components
	for _, sourceSettings := range config.Sources {
		name, ok := sourceSettings["name"]
		if !ok {
			logMissingField("Source", "name")
		}

		_, exists := sourceLookup[name]
		if exists {
			log.Fatalf("Duplicate source name in config: %s", name)
		}

		stype, ok := sourceSettings["type"]
		if !ok {
			logMissingField("Source", "type")
		}

		source := NewSource(stype, sourceSettings)
		sourceLookup[name] = source
	}

	for _, sinkSettings := range config.Sinks {
		name, ok := sinkSettings["name"]
		if !ok {
			logMissingField("Sink", "name")
		}

		_, exists := sinkLookup[name]
		if exists {
			log.Fatalf("Duplicate sink name in config: %s", name)
		}

		stype, ok := sinkSettings["type"]
		if !ok {
			logMissingField("Sink", "type")
		}

		sink := NewSink(stype, sinkSettings)
		sinkLookup[name] = sink
	}

	for _, channelSettings := range config.Channels {
		name, ok := channelSettings["name"]
		if !ok {
			logMissingField("Channel", "name")
		}

		_, exists := channelLookup[name]
		if exists {
			log.Fatalf("Duplicate channel name in config: %s", name)
		}

		ctype, ok := channelSettings["type"]
		if !ok {
			logMissingField("Channel", "type")
		}

		channel := NewChannel(ctype, channelSettings)
		channelLookup[name] = channel
	}

	// set up bindings
	for _, sourceSettings := range config.Sources {
		name := sourceSettings["name"]
		channelNames, ok := sourceSettings["channel"]
		if !ok {
			logMissingField("Source", "channel")
		}

		channels := make([]Channel, 0)
		for _, channelName := range strings.Split(channelNames, ",") {
			channelName = strings.TrimSpace(channelName)
			channel, exists := channelLookup[channelName]
			if !exists {
				log.Fatalf("Config for source named %s has invalid channel %s", name, channelName)
			}
			channels = append(channels, channel)
		}

		source := sourceLookup[name]

		for _, channel := range channels {
			source.SetChannel(channel)
		}
	}

	for _, sinkSettings := range config.Sinks {
		name := sinkSettings["name"]
		channelName, ok := sinkSettings["channel"]
		if !ok {
			logMissingField("Sink", "channel")
		}

		channel, exists := channelLookup[channelName]
		if !exists {
			log.Fatalf("Config for sink named %s has invalid channel %s", name, channelName)
		}

		sink := sinkLookup[name]
		sink.SetChannel(channel)
	}

	// start the channels first
	for _, channel := range channelLookup {
		channel.Start()
	}

	for _, sink := range sinkLookup {
		sink.Start()
	}

	for _, source := range sourceLookup {
		source.Start()
	}

	go ConfigReloader()
}

func logMissingField(componentType string, field string) {
	log.Fatalf("%s missing %s field in config", componentType, field)
}

func ConfigReloader() {
	tick := time.Tick(time.Second * 10)
	for _ = range tick {
		// do reloads
	}
}
