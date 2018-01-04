package dgo2poc

// Data for WSOPResume.
// Resuming sessions is handled internally and this should not be touched by end-users.
type wsResume struct {
	Token     string `json:"token"`
	SessionID string `json:"session_id"`
	Seq       int    `json:"seq"`
}
