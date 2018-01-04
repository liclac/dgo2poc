package dgo2poc

import (
	"encoding/json"
)

// A websocket payload.
type wsPayload struct {
	OP   WSOP            `json:"op"`
	Data json.RawMessage `json:"d,omitempty"`

	// Only used for OP0 (Dispatch)
	Seq  int    `json:"s,omitempty"`
	Type string `json:"t,omitempty"`
}
