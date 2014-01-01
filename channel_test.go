package main

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
)

func makeDummyEvents(count int) []Event {
	events := []Event{}
	for i := 0; i < count; i++ {
		e := NewEvent()
		e.Body = []byte(fmt.Sprintf("Event %d", i))
		e.Headers["num"] = strconv.Itoa(i)
		events = append(events, e)
	}
	return events
}

func ChannelAddEventTest(c Channel, t *testing.T) {
	events := makeDummyEvents(1)
	err := c.AddEvent(events[0])
	if err != nil {
		t.Fatalf("Failed to add event: %s", err)
	}
}

func ChannelAddEventsTest(c Channel, t *testing.T) {
	events := makeDummyEvents(2)
	err := c.AddEvents(events)
	if err != nil {
		t.Fatalf("Failed to add events: %s", err)
	}
}

func ChannelGetOldestTest(c Channel, t *testing.T) {
	events := makeDummyEvents(2)
	err := c.AddEvents(events)
	if err != nil {
		t.Fatalf("Failed to add events: %s", err)
	}

	count, returnedEvents, err := c.GetOldest(1)
	if err != nil {
		t.Fatalf("Failed to get event: %s", err)
	}
	if count != 1 {
		t.Errorf("Supposed to return 1 event, instead got %d", count)
	}
	if !bytes.Equal(returnedEvents[0].Body, events[0].Body) {
		t.Errorf("Got wrong event back, Expected Body: %s, got Body: %s", events[0].Body, returnedEvents[0].Body)
	}

	count, returnedEvents, err = c.GetOldest(2)
	if err != nil {
		t.Fatalf("Failed to get event: %s", err)
	}
	if count != 2 {
		t.Errorf("Supposed to return 2 events, instead got %d", count)
	}
	if !bytes.Equal(returnedEvents[0].Body, events[0].Body) || !bytes.Equal(returnedEvents[1].Body, events[1].Body) {
		t.Errorf("Got events in wrong order")
	}
}

func ChannelGetAllTest(c Channel, t *testing.T) {
	events := makeDummyEvents(3)

	err := c.AddEvents(events)
	if err != nil {
		t.Fatalf("Failed to add events: %s", err)
	}

	count, returnedEvents, err := c.GetAll()
	if err != nil {
		t.Fatalf("Failed to get events: %s", err)
	}
	if count != 3 {
		t.Errorf("Supposed to return 3 events, instead got %d", count)
	}
	if !bytes.Equal(returnedEvents[0].Body, events[0].Body) ||
		!bytes.Equal(returnedEvents[1].Body, events[1].Body) ||
		!bytes.Equal(returnedEvents[2].Body, events[2].Body) {
		t.Errorf("Got events in wrong order")
	}
}

func ChannelConfirmGetTest(c Channel, t *testing.T) {
	events := makeDummyEvents(3)

	err := c.AddEvents(events)
	if err != nil {
		t.Fatalf("Failed to add events: %s", err)
	}

	_, _, err = c.GetOldest(1)
	if err != nil {
		t.Fatalf("failed to get events: %s", err)
	}

	c.ConfirmGet(1)
	_, returnedEvents, err := c.GetOldest(1)
	if err != nil {
		t.Fatalf("failed to get events: %s", err)
	}
	if !bytes.Equal(returnedEvents[0].Body, events[1].Body) {
		t.Errorf("got wrong event back, expected body: %s, got body: %s", events[1].Body, returnedEvents[0].Body)
	}
}

func ChannelStartTest(c Channel, t *testing.T) {
	err := c.Start()
	if err != nil {
		t.Errorf("Got error starting channel: %s", err)
	}
}
