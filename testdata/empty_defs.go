package testdata

//go:generate go2jsonc -type EmptyDefs -out empty_defs.jsonc

// EmptySubType define a struct with non-initialized fields.
type EmptySubType struct {
	A string // Field A
	B int    // Field B
}

// EmptyDefs define a struct with non-initialized fields.
type EmptyDefs struct {
	Test1 EmptySubType
	Test2 []EmptySubType
}

func EmptyDefsDefaults() *EmptyDefs {
	return &EmptyDefs{}
}
