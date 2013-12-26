package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"sync"
	"time"
)

func init() {
	RegisterSink("legacy", NewLegacyFileSink)
}

type LegacyFileSink struct {
	channel        Channel
	transPerFile   uint
	rollPeriod     time.Duration
	incompletePath string
	completePath   string
	currentFile    *os.File
	fileLock       sync.Mutex
}

func NewLegacyFileSink(config ComponentSettings) Sink {
	incompletePath, ok := config["incomplete"]
	if !ok {
		log.Fatal("Must configure incomplete path for legacy sink")
	}

	completePath, ok := config["complete"]
	if !ok {
		log.Fatal("Must configure complete path for legacy sink")
	}

	startingFile, err := os.Create(path.Join(incompletePath, fmt.Sprintf("%d.txt.inc", time.Now().UTC().Unix())))
	if err != nil {
		log.Fatalf("legacysink: error opening file for writing: %s", err)
	}

	return &LegacyFileSink{
		channel:        nil,
		transPerFile:   uint(100),
		rollPeriod:     10 * time.Second,
		incompletePath: incompletePath,
		completePath:   completePath,
		currentFile:    startingFile}
}

func (l *LegacyFileSink) SetChannel(channel Channel) error {
	l.channel = channel
	return nil
}

func (l *LegacyFileSink) Start() error {
	go l.loopForever()
	return nil
}

func (l *LegacyFileSink) loopForever() {
	ticker := time.Tick(l.rollPeriod)
	txCount := uint(0)

	for {
		select {
		case <-ticker:
			if txCount == uint(0) {
				l.rollFile(true)
				txCount = 0
			} else {
				l.rollFile(false)
				txCount = 0
			}
		default:
			if l.channel == nil {
				continue
			}
			count, events, err := l.channel.GetAll()
			if err != nil {
				log.Printf("legacysink: channel get all: %s", err)
				continue
			}

			for _, event := range events {
				l.writeEvent(event)
				txCount += 1
				if txCount == l.transPerFile {
					txCount = 0
					l.rollFile(false)
				}
			}
			l.channel.ConfirmGet(count)
		}
	}
}

func (l *LegacyFileSink) rollFile(deleteOld bool) {
	l.fileLock.Lock()
	defer l.fileLock.Unlock()
	err := l.currentFile.Close()
	if err != nil {
		log.Fatalf("legacysink: close file for roll: %s", err)
	}

	oldName := l.currentFile.Name()
	newName := path.Base(oldName)
	newName = newName[:len(newName)-4]
	newName = path.Join(l.completePath, newName)

	if deleteOld {
		err = os.Remove(oldName)
		if err != nil {
			log.Fatalf("legacysink: remove file: %s")
		}
	} else {
		err = os.Rename(oldName, newName)
		if err != nil {
			log.Fatalf("legacysink: rename file: %s", err)
		}
	}
	l.currentFile, err = os.Create(l.getNewIncFilename())
	if err != nil {
		log.Fatalf("legacysink: open new incomplete: %s", err)
	}
}

func (l *LegacyFileSink) getNewIncFilename() string {
	return path.Join(l.incompletePath, fmt.Sprintf("%d.txt.inc", time.Now().UTC().Unix()))
}

func (l *LegacyFileSink) writeEvent(event Event) {
	msgOut := fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n",
		event.Headers["Timestamp"],
		event.Headers["RemoteAddr"],
		string(event.Body),
		event.Headers["UserAgent"],
		event.Headers["Referrer"])
	l.fileLock.Lock()
	defer l.fileLock.Unlock()
	fmt.Fprint(l.currentFile, msgOut)
}

func (l *LegacyFileSink) ReloadConfig(config ComponentSettings) bool {
	return true
}
