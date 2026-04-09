package types

import (
	"encoding/json"
)

type EventType string

var (
	MatchJoin  EventType = "MATCH_JOIN"
	MatchLeave EventType = "MATCH_LEAVE"
)

type EventEnvelop struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload,omitempty"`
}
