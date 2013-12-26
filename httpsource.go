package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

func init() {
	RegisterSource("http", NewHttpSource)
}

type HttpSource struct {
	channels []Channel
	server   *http.Server
}

func NewHttpSource(config ComponentSettings) Source {

	port, ok := config["port"]
	if !ok {
		log.Fatalf("Must configure port for http source")
	}

	path, ok := config["path"]
	if !ok {
		log.Fatalf("Must configure path for http source")
	}

	h := &HttpSource{channels: make([]Channel, 0)}

	mux := http.NewServeMux()
	mux.Handle(path, h)

	h.server = &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Starting http sink at http://localhost:%s%s", port, path)
	return h
}

func (h *HttpSource) SetChannel(c Channel) error {
	h.channels = append(h.channels, c)
	return nil
}

func (h *HttpSource) Start() error {
	go h.server.ListenAndServe()
	return nil
}

func (h *HttpSource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ts := time.Now().UTC().Unix()
	var body []byte
	var err error
	switch r.Method {
	case "GET":
		body = []byte(r.URL.RawQuery)
	case "POST":
		body, err = ioutil.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %s", err)
			return
		}
	default:
		log.Printf("Unsupported method: %s", r.Method)
		return
	}
	m := NewEvent()
	m.Body = body
	m.Headers["Timestamp"] = strconv.FormatInt(ts, 10)
	m.Headers["Referrer"] = r.Referer()
	m.Headers["UserAgent"] = r.UserAgent()
	m.Headers["RemoteAddr"] = r.RemoteAddr
	for _, channel := range h.channels {
		channel.AddEvent(m)
	}
}

func (h *HttpSource) ReloadConfig(config ComponentSettings) bool {
	return true
}
