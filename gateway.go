package dgo2poc

type Gateway struct {
	URL    string `json:"url"`
	Shards int    `json:"shards,omitempty"`
}
