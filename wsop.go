package dgo2poc

// OPCodes for websocket connections.
type WSOP int

const (
	// Dispatches an event.
	WSOPDispatch WSOP = iota
	// Used for ping checking.
	WSOPHeartbeat
	// Used for client handshake.
	WSOPIdentify
	// Used to update the client status.
	WSOPStatusUpdate
	// Used to join/move/leave voice channels.
	WSOPVoiceStateUpdate
	// Used for voice ping checking.
	WSOPVoiceServerPing
	// Used to resume a closed connection.
	WSOPResume
	// Used to tell clients to reconnect to the gateway.
	WSOPReconnect
	// Used to request guild members.
	WSOPRequestGuildMembers
	// Used to notify client they have an invalid session id.
	WSOPInvalidSession
	// Sent immediately after connecting, contains heartbeat and server debug information.
	WSOPHello
	// Sent immediately following a client heartbeat that was received.
	WSOPHeartbeatAck
)
