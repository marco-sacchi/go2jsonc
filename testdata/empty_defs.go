package testdata

//go:generate go2jsonc -type EmptyDefs -out empty_defs.jsonc
//go:generate go2jsonc -type EmptyDefs -doc-types NotFields -out empty_defs_fields.jsonc
//go:generate go2jsonc -type EmptyDefs -doc-types NotStructFields -out empty_defs_struct.jsonc
//go:generate go2jsonc -type EmptyDefs -doc-types NotArrayFields -out empty_defs_array.jsonc
//go:generate go2jsonc -type EmptyDefs -doc-types NotMapFields -out empty_defs_map.jsonc

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
