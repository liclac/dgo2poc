package dgo2poc

// Data for WSOPIdentify.
type wsIdentify struct {
	Token          string          `json:"token"`
	Properties     wsIdentifyProps `json:"properties"`
	Compress       bool            `json:"compress"`
	LargeThreshold int             `json:"large_threshold"`
	Shard          [2]int          `json:"shard"`
	Presence       WSStatus        `json:"presence"`
}

// Properties for an Identify payload.
type wsIdentifyProps struct {
	OS      string `json:"$os"`      // OS family.
	Browser string `json:"$browser"` // Library name.
	Device  string `json:"$device"`  // Library name.
}
