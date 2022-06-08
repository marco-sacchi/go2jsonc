package go2jsonc

import (
	"os"
	"strings"
	"testing"
)

func TestGenerator(t *testing.T) {
	var tests = []struct {
		typeName string
		filename string
	}{
		{"Embedding", "./testdata/embedding.jsonc"},
		{"Empty", "./testdata/empty.jsonc"},
		{"Nesting", "./testdata/nesting.jsonc"},
		{"Simple", "./testdata/simple.jsonc"},
	}

	whitespacesReplacer := strings.NewReplacer(" ", "◦", "\t", "➞➞➞➞")
	for _, test := range tests {
		jsonc, err := Generate("./testdata", test.typeName)
		if err != nil {
			t.Fatal(err)
		}

		content, err := os.ReadFile(test.filename)
		if err != nil {
			t.Fatal(err)
		}

		want := string(content)

		if jsonc != want {
			t.Fatalf("Generated JSONC mismatch %s:\n%s\n\nwant %s:\n%s",
				test.typeName,
				whitespacesReplacer.Replace(jsonc),
				test.filename,
				whitespacesReplacer.Replace(want))
		}
	}
}
