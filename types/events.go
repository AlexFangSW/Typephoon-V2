package types

import (
	"encoding/json"
)

type Event string

var Events = struct {
	MatchJoin  Event
	MatchLeave Event
}{
	MatchJoin:  "MATCH_JOIN",
	MatchLeave: "MATCH_LEAVE",
}

type Header string

var Headers = struct {
	UserID Header
	GameID Header
}{
	UserID: "User-ID",
	GameID: "Game-ID",
}

type EventEnvelop struct {
	Type    Event           `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}
