package testdata

//go:generate go2jsonc -type=Simple -out=simple.jsonc

// Simple defines a simple user.
type Simple struct {
	// User name documentation block.
	Name    string // User name comment.
	Surname string // User surname comment.

	// Age documentation block.
	Age        int `json:"age"`         // User age.
	StarsCount int `json:"stars_count"` // Number of stars achieved.

	Addresses []string // Addresses comment.
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
	}
}
