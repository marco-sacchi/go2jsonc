package testdata

//go:generate go2jsonc -type Embedding -out embedding.jsonc
//go:generate go2jsonc -type Embedding -doc-types NotFields -out embedding_not_fields.jsonc
//go:generate go2jsonc -type Embedding -doc-types NotStructFields -out embedding_not_struct.jsonc
//go:generate go2jsonc -type Embedding -doc-types NotArrayFields -out embedding_not_array.jsonc
//go:generate go2jsonc -type Embedding -doc-types NotMapFields -out embedding_not_map.jsonc

// Embedded test struct.
type Embedded struct {
	// Identifier documentation block.
	Identifier int  `json:"id"`
	Enabled    bool // Enabled comment line.

	Reserved uint32 `json:"reserved"`
}

// Embedding test struct.
type Embedding struct {
	// Embedded documentation block.
	Embedded

	Position float32 `json:"position"` // Position comment line.
	// Velocity documentation block.
	Velocity     float32 `json:"velocity"`
	Acceleration float32 `json:"accel"`

	Reserved string `json:"reserved"` // Shadowing field.
}

func EmbeddingDefaults() *Embedding {
	return &Embedding{
		Embedded: Embedded{
			Identifier: 1234,
			Enabled:    false,
			Reserved:   0x10,
		},
		Position:     1.0,
		Velocity:     2,
		Acceleration: 0.23,
		Reserved:     "Shadowing",
	}
}
