package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"sync"
)

func init() {
	RegisterChannel("sqlite", NewSqliteChannel)
}

type SqliteChannel struct {
	dbLock          sync.RWMutex
	db              *sql.DB
	unconfirmedGets []int
}

func NewSqliteChannel(config ComponentSettings) Channel {
	dbPath, ok := config["db"]
	if !ok {
		log.Fatal("must configure db for sqlite channel")
	}

	sqliteChannel := &SqliteChannel{}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	sqliteChannel.db = db
	err = sqliteChannel.initDb()
	if err != nil {
		log.Fatal(err)
	}
	sqliteChannel.unconfirmedGets = make([]int, 0)

	return sqliteChannel
}

func (s *SqliteChannel) initDb() error {
	sql := `
create table if not exists queue (
id integer primary key autoincrement,
body BLOB);`
	_, err := s.db.Exec(sql)
	if err != nil {
		log.Fatalf("%q: %s\n", err, sql)
	}
	return nil
}

func (s *SqliteChannel) AddEvent(m Event) error {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	encoded, err := json.Marshal(m)
	if err != nil {
		return err
	}
	_, err = s.db.Exec("insert into queue (body) values (?)", encoded)
	return err
}

func (s *SqliteChannel) AddEvents(m []Event) error {
	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	for _, event := range m {
		encoded, err := json.Marshal(event)
		if err != nil {
			return err
		}
		_, err = s.db.Exec("insert into queue (body) values (?)", encoded)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *SqliteChannel) GetOldest(count int) (int, []Event, error) {
	s.dbLock.RLock()
	defer s.dbLock.RUnlock()

	events := make([]Event, 0, count)
	ids := make([]int, 0, count)

	rows, err := s.db.Query("select id, body from queue order by id limit ?", count)
	if err != nil {
		return 0, []Event{}, err
	}
	for rows.Next() {
		var id int
		var encoded []byte
		if err := rows.Scan(&id, &encoded); err != nil {
			return 0, []Event{}, err
		}
		var m Event
		err = json.Unmarshal(encoded, &m)
		if err != nil {
			return 0, []Event{}, err
		}
		events = append(events, m)
		ids = append(ids, id)
	}
	s.unconfirmedGets = append(s.unconfirmedGets, ids...)
	return len(events), events, nil
}

func (s *SqliteChannel) GetAll() (int, []Event, error) {
	s.dbLock.RLock()
	defer s.dbLock.RUnlock()

	events := make([]Event, 0)
	ids := make([]int, 0)

	rows, err := s.db.Query("select id, body from queue order by id")
	if err != nil {
		return 0, []Event{}, err
	}

	for rows.Next() {
		var id int
		var encoded []byte
		if err := rows.Scan(&id, &encoded); err != nil {
			return 0, []Event{}, err
		}
		var m Event
		err = json.Unmarshal(encoded, &m)
		if err != nil {
			return 0, []Event{}, err
		}
		events = append(events, m)
		ids = append(ids, id)
	}
	s.unconfirmedGets = append(s.unconfirmedGets, ids...)
	return len(events), events, nil
}

func (s *SqliteChannel) ConfirmGet(count int) error {
	if count == 0 {
		s.unconfirmedGets = nil
		return nil
	}

	s.dbLock.Lock()
	defer s.dbLock.Unlock()

	ix := IntMin(count-1, len(s.unconfirmedGets)-1)
	max_id := s.unconfirmedGets[ix]
	_, err := s.db.Exec("delete from queue where id <= ?", max_id)
	if err != nil {
		return err
	}

	s.unconfirmedGets = nil
	return nil
}

func (s *SqliteChannel) Start() error {
	return nil
}

func (s *SqliteChannel) ReloadConfig(config ComponentSettings) bool {
	return true
}
