package types

import (
	"encoding/json"
)

type Subject string

var Subjects = struct {
	MatchJoin Subject
}{
	MatchJoin: "match.join",
}

type Event string

var Events = struct {
	MatchJoin  Event
	MatchLeave Event
	Error      Event
}{
	MatchJoin:  "MATCH_JOIN",
	MatchLeave: "MATCH_LEAVE",
	Error:      "Error",
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

type ErrorPayload struct {
	Error string `json:"error"`
}
