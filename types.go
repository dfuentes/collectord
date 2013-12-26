package main

import (
	"log"
)

type Event struct {
	Headers map[string]string
	Body    []byte
}

type Channel interface {
	// Supporting Sources
	AddEvent(Event) error
	AddEvents([]Event) error

	// Supporting Sinks
	GetOldest(int) (int, []Event, error)
	GetAll() (int, []Event, error)
	ConfirmGet(int) error

	Start() error

	ReloadConfig(config ComponentSettings) bool
}

type Sink interface {
	SetChannel(Channel) error
	Start() error

	ReloadConfig(config ComponentSettings) bool
}

type Source interface {
	SetChannel(Channel) error
	Start() error

	ReloadConfig(config ComponentSettings) bool
}

func NewEvent() Event {
	return Event{make(map[string]string), make([]byte, 0)}
}

// Global source registry

var registeredSources map[string]func(ComponentSettings) Source = make(map[string]func(ComponentSettings) Source)

func RegisterSource(name string, constructor func(ComponentSettings) Source) {
	registeredSources[name] = constructor
}

func NewSource(name string, config ComponentSettings) Source {
	constructor, ok := registeredSources[name]
	if !ok {
		log.Fatalf("No source registered for name [%s]", name)
	}
	return constructor(config)
}

// Global channel registry

var registeredChannels map[string]func(ComponentSettings) Channel = make(map[string]func(ComponentSettings) Channel)

func RegisterChannel(name string, constructor func(ComponentSettings) Channel) {
	registeredChannels[name] = constructor
}

func NewChannel(name string, config ComponentSettings) Channel {
	constructor, ok := registeredChannels[name]
	if !ok {
		log.Fatalf("No channel registered for name [%s]", name)
	}
	return constructor(config)
}

// Global sink registry

var registeredSinks map[string]func(ComponentSettings) Sink = make(map[string]func(ComponentSettings) Sink)

func RegisterSink(name string, constructor func(ComponentSettings) Sink) {
	registeredSinks[name] = constructor
}

func NewSink(name string, config ComponentSettings) Sink {
	constructor, ok := registeredSinks[name]
	if !ok {
		log.Fatalf("No sink registered for name [%s]", name)
	}
	return constructor(config)
}
