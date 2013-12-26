package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"testing"
)

// TODO: test error conditions

func getInitConfig() ComponentSettings {
	temp := os.TempDir()
	db := path.Join(temp, "collect_test_db")
	c := ComponentSettings{"db": db}
	return c
}

func cleanupDb(config ComponentSettings) {
	db, ok := config["db"]
	if !ok {
		return
	}

	err := os.Remove(db)
	if err != nil {
		log.Printf("error cleaning up db: %s", err)
	}
}

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

func TestSqliteChannelInit(t *testing.T) {
	c := getInitConfig()
	defer cleanupDb(c)

	// will cause exit if init fails
	sqliteChannel := NewSqliteChannel(c)
	defer sqliteChannel.(*SqliteChannel).db.Close()
}

func TestSqliteChannelAddEvent(t *testing.T) {
	c := getInitConfig()
	defer cleanupDb(c)

	sqliteChannel := NewSqliteChannel(c)
	defer sqliteChannel.(*SqliteChannel).db.Close()

	events := makeDummyEvents(1)

	err := sqliteChannel.AddEvent(events[0])
	if err != nil {
		t.Fatalf("Failed to add event: %s", err)
	}
}

func TestSqliteChannelAddEvents(t *testing.T) {
	c := getInitConfig()
	defer cleanupDb(c)

	sqliteChannel := NewSqliteChannel(c)
	defer sqliteChannel.(*SqliteChannel).db.Close()

	events := makeDummyEvents(2)

	err := sqliteChannel.AddEvents(events)
	if err != nil {
		t.Fatalf("Failed to add events: %s", err)
	}
}

func TestSqliteChannelGetOldest(t *testing.T) {
	c := getInitConfig()
	defer cleanupDb(c)

	sqliteChannel := NewSqliteChannel(c)
	defer sqliteChannel.(*SqliteChannel).db.Close()

	events := makeDummyEvents(2)

	err := sqliteChannel.AddEvents(events)
	if err != nil {
		t.Fatalf("Failed to add events: %s", err)
	}

	count, returnedEvents, err := sqliteChannel.GetOldest(1)
	if err != nil {
		t.Fatalf("Failed to get event: %s", err)
	}
	if count != 1 {
		t.Errorf("Supposed to return 1 event, instead got %d", count)
	}

	if !bytes.Equal(returnedEvents[0].Body, events[0].Body) {
		t.Errorf("Got wrong event back, Expected Body: %s, got Body: %s", events[0].Body, returnedEvents[0].Body)
	}

	count, returnedEvents, err = sqliteChannel.GetOldest(2)
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

func TestSqliteChannelGetAll(t *testing.T) {
	c := getInitConfig()
	defer cleanupDb(c)

	sqliteChannel := NewSqliteChannel(c)
	defer sqliteChannel.(*SqliteChannel).db.Close()

	events := makeDummyEvents(3)

	err := sqliteChannel.AddEvents(events)
	if err != nil {
		t.Fatalf("Failed to add events: %s", err)
	}

	count, returnedEvents, err := sqliteChannel.GetAll()
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

func TestSqliteChannelConfirmGet(t *testing.T) {
	c := getInitConfig()
	defer cleanupDb(c)

	sqliteChannel := NewSqliteChannel(c)
	defer sqliteChannel.(*SqliteChannel).db.Close()

	events := makeDummyEvents(3)

	err := sqliteChannel.AddEvents(events)
	if err != nil {
		t.Fatalf("Failed to add events: %s", err)
	}

	_, _, err = sqliteChannel.GetOldest(1)
	if err != nil {
		t.Fatalf("failed to get events: %s", err)
	}

	sqliteChannel.ConfirmGet(1)
	_, returnedEvents, err := sqliteChannel.GetOldest(1)
	if err != nil {
		t.Fatalf("failed to get events: %s", err)
	}
	if !bytes.Equal(returnedEvents[0].Body, events[1].Body) {
		t.Errorf("got wrong event back, expected body: %s, got body: %s", events[1].Body, returnedEvents[0].Body)
	}
}

func TestSqliteChannelStart(t *testing.T) {
	c := getInitConfig()
	defer cleanupDb(c)

	sqliteChannel := NewSqliteChannel(c)
	defer sqliteChannel.(*SqliteChannel).db.Close()

	err := sqliteChannel.Start()
	if err != nil {
		t.Errorf("Got error starting channel: %s", err)
	}
}
