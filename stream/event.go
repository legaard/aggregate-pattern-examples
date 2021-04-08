package stream

import "time"

type Metadata struct {
	StreamID    string
	StreamType  string
	EventNumber int
	EventTime   time.Time
	EventID     string
}

type Event struct {
	Metadata
	Data interface{}
}
