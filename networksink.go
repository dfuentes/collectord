package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"time"
)

func init() {
	RegisterSink("gob", NewGobSink)
}

type GobSink struct {
	channel Channel
	encBuf  *bufio.Writer
	enc     *gob.Encoder
	conn    net.Conn
	host    string
	port    string
}

func NewGobSink(config ComponentSettings) Sink {
	host, ok := config["host"]
	if !ok {
		log.Fatal("must configure host for gob sink")
	}

	port, ok := config["port"]
	if !ok {
		log.Fatal("must configure port for gob sink")
	}

	gs := &GobSink{}
	gs.host = host
	gs.port = fmt.Sprintf(":%s", port)

	//	gs.setupConnection()

	return gs
}

func (gs *GobSink) setupConnection() error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s%s", gs.host, gs.port))
	if err != nil {
		log.Printf("Unable to connect to server at %s%s: %s", gs.host, gs.port, err)
		gs.conn = nil
		gs.encBuf = nil
		gs.enc = nil
		return err
	}
	log.Printf("Connected to server @ %s%s", gs.host, gs.port)
	gs.conn = conn
	gs.conn.(*net.TCPConn).SetWriteBuffer(4096)
	gs.encBuf = bufio.NewWriter(conn)
	gs.enc = gob.NewEncoder(gs.encBuf)
	return nil
}

func (gs *GobSink) Start() error {
	go gs.loopForever()
	return nil
}

func (gs *GobSink) loopForever() {
	tick := time.Tick(time.Millisecond * 500)

	//	lastFailedConnection := time.Time{}
	attempt := 0

	err := gs.setupConnection()
	if err != nil {
		//		lastFailedConnection = time.Now()
		attempt = 1
	}

mainfor:
	for _ = range tick {
		if gs.channel == nil {
			continue
		}
		if gs.conn == nil {
			err = gs.setupConnection()
			if err != nil {
				//				lastFailedConnection = time.Now()
				attempt = attempt + 1
				log.Printf("gobsink: Failed to connect, retries: %d", attempt)
				continue
			} else {
				attempt = 0
			}
		}

		count, events, err := gs.channel.GetAll()
		if err != nil {
			log.Printf("gobsink: Error getting events from channel: %s", err)
			continue
		}

		if shouldSendDummy(events) {
			// send dummy event if we are below threshold
			// helps to find broken connections
			if err = gs.enc.Encode(Event{}); err != nil {
				log.Printf("gobsink: dummy enc: %s", err)
				gs.abortSend()
				continue
			}
			if err = gs.encBuf.Flush(); err != nil {
				log.Printf("gobsink: dummy send: %s", err)
				gs.abortSend()
				continue
			}
		}

		for _, event := range events {
			if err = gs.enc.Encode(event); err != nil {
				log.Printf("gobsink: encode: %s", err)
				gs.abortSend()
				continue mainfor
			}
		}
		err = gs.encBuf.Flush()
		if err != nil {
			log.Printf("gobsink: Error flushing encoding buffer: %s", err)
			gs.abortSend()
			continue
		}
		gs.channel.ConfirmGet(count)
	}
}

func (gs *GobSink) abortSend() {
	gs.conn.Close()
	gs.conn = nil
	gs.channel.ConfirmGet(0)
}

func (gs *GobSink) SetChannel(channel Channel) error {
	gs.channel = channel
	return nil
}

func shouldSendDummy(m []Event) bool {
	if len(m) == 0 {
		return false
	}

	size := 0
	for _, event := range m {
		size += len(event.Body)
	}

	if size < MIN_PACKET_THRESHOLD {
		return true
	}

	return false
}

func (gs *GobSink) ReloadConfig(config ComponentSettings) bool {
	return true
}
