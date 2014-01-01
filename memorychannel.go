package main

import (
	"container/list"
	"sync"
)

func init() {
	RegisterChannel("memory", NewMemoryChannel)
}

type MemoryChannel struct {
	queue *list.List
	lock  sync.Mutex
}

func NewMemoryChannel(config ComponentSettings) Channel {
	return &MemoryChannel{
		queue: list.New(),
	}
}

func (m *MemoryChannel) AddEvent(e Event) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.queue.PushFront(e)
	return nil
}

func (m *MemoryChannel) AddEvents(e []Event) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, event := range e {
		m.queue.PushFront(event)
	}
	return nil
}

func (m *MemoryChannel) GetOldest(count int) (int, []Event, error) {
	numToGet := IntMin(m.queue.Len(), count)
	events := make([]Event, 0, numToGet)
	back := m.queue.Back()
	for i := 0; i < numToGet; i++ {
		events = append(events, back.Value.(Event))
		back = back.Prev()
	}
	return numToGet, events, nil
}

func (m *MemoryChannel) GetAll() (int, []Event, error) {
	events := make([]Event, 0, m.queue.Len())
	for e := m.queue.Back(); e != nil; e = e.Prev() {
		events = append(events, e.Value.(Event))
	}
	return len(events), events, nil
}

func (m *MemoryChannel) ConfirmGet(count int) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	numToConfirm := IntMin(m.queue.Len(), count)
	back := m.queue.Back()
	for i := 0; i < numToConfirm; i++ {
		n := back.Prev()
		m.queue.Remove(back)
		back = n
	}
	return nil
}

func (m *MemoryChannel) Start() error {
	return nil
}

func (m *MemoryChannel) ReloadConfig(config ComponentSettings) bool {
	return true
}
