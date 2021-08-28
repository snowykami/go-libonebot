package libonebot

import (
	"encoding/json"
	"sync"
	"time"
)

type eventType struct{ string }

func (t eventType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.string)
}

var (
	EventTypeMessage eventType = eventType{"message"}
	EventTypeNotice  eventType = eventType{"notice"}
	EventTypeRequest eventType = eventType{"request"}
	EventTypeMeta    eventType = eventType{"meta"}
)

type Event struct {
	lock       sync.Mutex
	Platform   string    `json:"platform"`
	Time       int64     `json:"time"`
	SelfID     string    `json:"self_id"`
	Type       eventType `json:"type"`
	DetailType string    `json:"detail_type"`
}

type AnyEvent interface {
	TryFixUp() bool
}

func (e *Event) TryFixUp() bool {
	e.lock.Lock()
	defer e.lock.Unlock()

	if e.Platform == "" || e.SelfID == "" || e.Type.string == "" || e.DetailType == "" {
		return false
	}
	if e.Time == 0 {
		e.Time = time.Now().Unix()
	}
	return true
}

type MessageEvent struct {
	Event
	UserID  string  `json:"user_id"`
	GroupID string  `json:"group_id,omitempty"`
	Message Message `json:"message"`
}

type NoticeEvent struct {
	Event
}

type RequestEvent struct {
	Event
}

type MetaEvent struct {
	Event
}