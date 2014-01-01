package main

import (
	"log"
	"os"
	"path"
	"testing"
)

// TODO: test error conditions

func initSqliteChannelTest() (ComponentSettings, Channel) {
	temp := os.TempDir()
	db := path.Join(temp, "collect_test_db")
	c := ComponentSettings{"db": db}
	sqliteChannel := NewSqliteChannel(c)
	return c, sqliteChannel
}

func cleanupSqliteChannelTest(config ComponentSettings, channel Channel) {
	channel.(*SqliteChannel).db.Close()

	db, ok := config["db"]
	if !ok {
		return
	}

	err := os.Remove(db)
	if err != nil {
		log.Printf("error cleaning up db: %s", err)
	}
}

func TestSqliteChannelAddEvent(t *testing.T) {
	c, sqliteChannel := initSqliteChannelTest()
	defer cleanupSqliteChannelTest(c, sqliteChannel)

	ChannelAddEventTest(sqliteChannel, t)
}

func TestSqliteChannelAddEvents(t *testing.T) {
	c, sqliteChannel := initSqliteChannelTest()
	defer cleanupSqliteChannelTest(c, sqliteChannel)

	ChannelAddEventsTest(sqliteChannel, t)
}

func TestSqliteChannelGetOldest(t *testing.T) {
	c, sqliteChannel := initSqliteChannelTest()
	defer cleanupSqliteChannelTest(c, sqliteChannel)

	ChannelGetOldestTest(sqliteChannel, t)
}

func TestSqliteChannelGetAll(t *testing.T) {
	c, sqliteChannel := initSqliteChannelTest()
	defer cleanupSqliteChannelTest(c, sqliteChannel)

	ChannelGetAllTest(sqliteChannel, t)
}

func TestSqliteChannelConfirmGet(t *testing.T) {
	c, sqliteChannel := initSqliteChannelTest()
	defer cleanupSqliteChannelTest(c, sqliteChannel)

	ChannelConfirmGetTest(sqliteChannel, t)
}

func TestSqliteChannelStart(t *testing.T) {
	c, sqliteChannel := initSqliteChannelTest()
	defer cleanupSqliteChannelTest(c, sqliteChannel)

	ChannelStartTest(sqliteChannel, t)
}
