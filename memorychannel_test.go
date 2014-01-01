package main

import (
	"testing"
)

func initMemoryChannelTest() (ComponentSettings, Channel) {
	c := ComponentSettings{}
	memoryChannel := NewMemoryChannel(c)
	return c, memoryChannel
}

func cleanupMemoryChannelTest(config ComponentSettings, channel Channel) {
}

func TestMemoryChannelAddEvent(t *testing.T) {
	c, memoryChannel := initMemoryChannelTest()
	defer cleanupMemoryChannelTest(c, memoryChannel)

	ChannelAddEventTest(memoryChannel, t)
}

func TestMemoryChannelAddEvents(t *testing.T) {
	c, memoryChannel := initMemoryChannelTest()
	defer cleanupMemoryChannelTest(c, memoryChannel)

	ChannelAddEventsTest(memoryChannel, t)
}

func TestMemoryChannelGetOldest(t *testing.T) {
	c, memoryChannel := initMemoryChannelTest()
	defer cleanupMemoryChannelTest(c, memoryChannel)

	ChannelGetOldestTest(memoryChannel, t)
}

func TestMemoryChannelGetAll(t *testing.T) {
	c, memoryChannel := initMemoryChannelTest()
	defer cleanupMemoryChannelTest(c, memoryChannel)

	ChannelGetAllTest(memoryChannel, t)
}

func TestMemoryChannelConfirmGet(t *testing.T) {
	c, memoryChannel := initMemoryChannelTest()
	defer cleanupMemoryChannelTest(c, memoryChannel)

	ChannelConfirmGetTest(memoryChannel, t)
}

func TestMemoryChannelStart(t *testing.T) {
	c, memoryChannel := initMemoryChannelTest()
	defer cleanupMemoryChannelTest(c, memoryChannel)

	ChannelStartTest(memoryChannel, t)
}
