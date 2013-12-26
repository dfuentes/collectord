package main

import (
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
)

func init() {
	RegisterSource("gob", NewGobSource)
}

type GobSource struct {
	channels []Channel
	port     string
}

func NewGobSource(config ComponentSettings) Source {
	port, ok := config["port"]
	if !ok {
		log.Fatal("must set port for gob source")
	}

	return &GobSource{port: fmt.Sprintf(":%s", port), channels: make([]Channel, 0)}
}

func (g *GobSource) SetChannel(channel Channel) error {
	g.channels = append(g.channels, channel)
	return nil
}

func (g *GobSource) Start() error {
	go g.serveForever()
	return nil
}

func (g *GobSource) serveForever() {
	ln, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0%s", g.port))
	if err != nil {
		log.Fatalf("gobsource: failed to listen on port %s: %s", g.port, err)
	}
	defer ln.Close()

	log.Printf("gobsource: listening on port %s", g.port)

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("gobsource: failed to accept connection: %s", err)
			continue
		}
		log.Printf("gobsource: received connection from %s", conn.RemoteAddr())
		go g.handleConn(conn)
	}
}

func (g *GobSource) handleConn(conn net.Conn) {
	conn.(*net.TCPConn).SetLinger(0)
	decoder := gob.NewDecoder(conn)
	for {
		m := Event{}
		err := decoder.Decode(&m)

		if err == io.EOF {
			log.Printf("gobsource: connection closed by remote client %s", conn.RemoteAddr())
			conn.Close()
			return
		}
		for _, channel := range g.channels {
			channel.AddEvent(m)
		}
	}

}

func (g *GobSource) ReloadConfig(config ComponentSettings) bool {
	return true
}
