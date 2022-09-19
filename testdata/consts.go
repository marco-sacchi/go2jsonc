package testdata

type ConstType int

const (
	// ConstTypeA doc block.
	ConstTypeA ConstType = iota // ConstTypeA comment.
	ConstTypeB                  // ConstTypeB comment.
	// ConstTypeC doc block.
	ConstTypeC // ConstTypeC comment.
)

const (
	// ConstTypeD doc block.
	ConstTypeD ConstType = 1 << (iota + 5)
	// ConstTypeE doc block.
	ConstTypeE // ConstTypeE comment.
	// ConstTypeF doc block.
	ConstTypeF // ConstTypeF comment.
)

// ConstTypeEnum enumeration of const types.
var ConstTypeEnum = [ConstTypeC + 1]ConstType{
	ConstTypeA,
	ConstTypeB,
	ConstTypeC,
}

// ConstTypeABC lists A, B, C const types.
var ConstTypeABC = []ConstType{ConstTypeA, ConstTypeB, ConstTypeC}

// ConstTypeString is a constant of const types names.
var ConstTypeString = [ConstTypeC + 1]string{
	"ConstTypeA", "ConstTypeB", "ConstTypeC",
}
