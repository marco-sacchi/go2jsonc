package stats

// Info reports statistical info.
type Info struct {
	// PacketLoss documentation block.
	PacketLoss    int `json:"packet_loss"`     // Packet loss comment.
	RoundTripTime int `json:"round_trip_time"` // Round-trip time in milliseconds.
}
