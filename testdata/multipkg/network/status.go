package network

type ConnState int

const (
	// StateDisconnected signals the Disconnected state.
	StateDisconnected ConnState = iota
	// StateConnecting signals the connection-pending state.
	StateConnecting
	// StateConnected signals the Connected state.
	StateConnected
)

const (
	// StateFailed signals the Failed state.
	StateFailed ConnState = iota + 5
	// StateReconnecting signals the Reconnecting state.
	StateReconnecting
)

// Status reports connection status.
type Status struct {
	Connected bool      // Connected flag comment.
	State     ConnState // Connection state comment.
}
