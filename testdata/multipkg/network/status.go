package network

type ConnState int

const (
	// StateDisconnected signals the Disconnected state.
	StateDisconnected ConnState = iota // StateDisconnected comment.
	StateConnecting                    // StateConnecting comment.
	// StateConnected signals the Connected state.
	StateConnected // StateConnected comment.
)

const (
	// StateFailed signals the Failed state.
	StateFailed ConnState = iota + 5 // StateFailed comment.
	// StateReconnecting signals the Reconnecting state.
	StateReconnecting // StateReconnecting comment.
)

// Status reports connection status.
type Status struct {
	Connected bool      // Connected flag comment.
	State     ConnState // Connection state comment.
}
