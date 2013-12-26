package main

import (
	"log"
	"time"
)

func init() {
	RegisterSink("console", func(ComponentSettings) Sink { return &ConsoleSink{} })
}

type ConsoleSink struct {
	channel Channel
}

func (c *ConsoleSink) SetChannel(ch Channel) error {
	c.channel = ch
	return nil
}

func (c *ConsoleSink) Start() error {
	go c.loopForever()
	return nil
}

func (c *ConsoleSink) loopForever() {
	tick := time.Tick(time.Millisecond * 500)
	for _ = range tick {
		if c.channel == nil {
			continue
		}
		count, events, err := c.channel.GetAll()
		if err != nil {
			log.Printf("Error getting eventss: %s", err)
		}
		for _, event := range events {
			log.Printf("headers: %+v body: %s", event.Headers, event.Body)
		}
		c.channel.ConfirmGet(count)
	}
}

func (c *ConsoleSink) ReloadConfig(config ComponentSettings) bool {
	return true
}
