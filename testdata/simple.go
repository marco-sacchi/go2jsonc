package testdata

//go:generate go2jsonc -type Simple -out simple.jsonc
//go:generate go2jsonc -type Simple -doc-types NotFields -out simple_not_fields.jsonc
//go:generate go2jsonc -type Simple -doc-types NotStructFields -out simple_not_struct.jsonc
//go:generate go2jsonc -type Simple -doc-types NotArrayFields -out simple_not_array.jsonc
//go:generate go2jsonc -type Simple -doc-types NotMapFields -out simple_not_map.jsonc

// Simple defines a simple user.
type Simple struct {
	// User name documentation block.
	Name    string // User name comment.
	Surname string // User surname comment.

	// Age documentation block.
	Age        int `json:"age"`         // User age.
	StarsCount int `json:"stars_count"` // Number of stars achieved.

	Addresses []string // Addresses comment.

	Tags map[string]string // User tags.

	// Type documentation block.
	Type ConstType // Type of constant.
}

func SimpleDefaults() *Simple {
	return &Simple{
		Name:       "John",
		Surname:    "Doe",
		Age:        30,
		StarsCount: 5,
		Addresses: []string{
			"Address 1",
			"Address 2",
			"Address 3",
		},
		Tags: map[string]string{
			"Key1": "Value1",
			"Key2": "Value2",
			"Key3": "Value3",
		},
	}
}
