package testdata

//go:generate go2jsonc -type=Nesting -out nesting.jsonc

// Protocol defines a network protocol and version.
type Protocol struct {
	// Name describes the protocol name.
	// Multiple line documentation test.
	Name string // Protocol name.

	Major int // Major version.
	Minor int // Minor version.
}

// Nesting checks for correct struct nesting.
type Nesting struct {
	IP   string // Remote IP address.
	Port int    // Remote port.

	Default   Protocol   `json:"default_proto"`   // Default protocol.
	Optionals []Protocol `json:"optional_protos"` // Optional supported protocols.
}

func NestingDefaults() *Nesting {
	return &Nesting{
		IP:   "127.0.0.1",
		Port: 12345,
		Default: Protocol{
			Name:  "TCP",
			Major: 1,
			Minor: 0,
		},
		Optionals: []Protocol{
			{
				Name:  "UDP",
				Major: 1,
				Minor: 0,
			},
			{
				Name:  "HTTP",
				Major: 1,
				Minor: 1,
			},
		},
	}
}
