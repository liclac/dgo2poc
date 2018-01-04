package dgo2poc

type WSStatusType string

const (
	WSStatusOnline    = "online"
	WSStatusDND       = "dnd"
	WSStatusIdle      = "idle"
	WSStatusInvisible = "invisible"
	WSStatusOffline   = "offline"
)

type WSStatusGameType int

const (
	WSStatusGameIsGame WSStatusGameType = iota
	WSStatusGameIsStreaming
	WSStatusGameIsListening
	WSStatusGameIsWatching
)

// Data for WSOPStatusUpdate.
type WSStatus struct {
	Since  *int64        `json:"since"`
	Game   *WSStatusGame `json:"game"`
	Status string        `json:"status"`
	AFK    bool          `json:"afk"`
}

// Game for a WSStatus payload.
type WSStatusGame struct {
	Name string           `json:"name"`
	Type WSStatusGameType `json:"type"`

	// Stream URL, only used for WSStatusIsStreaming.
	URL string `json:"url,omitempty"`
}
