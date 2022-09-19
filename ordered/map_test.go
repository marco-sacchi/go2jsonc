package ordered

import "testing"

func TestMap(t *testing.T) {
	m := NewMap()

	if m.Len() != 0 {
		t.Fatalf("Wrong empty map length: got %d, want 0", m.Len())
	}

	items := []struct {
		Key   string
		Value string
		Index int
	}{
		{Key: "0", Value: "a", Index: 0},
		{Key: "1", Value: "b", Index: 1},
		{Key: "2", Value: "c", Index: 2},
		{Key: "3", Value: "d", Index: 3},
	}

	for _, item := range items {
		m.Append(item.Key, item.Value)
	}

	if m.Len() != len(items) {
		t.Fatalf("Wrong map length: got %d, want %d", m.Len(), len(items))
	}

	s := m.String()
	wantString := "ordered.Map[0:a 1:b 2:c 3:d]"
	if s != wantString {
		t.Fatalf("Wrong map String result: got:\n%s\nwant:\n%s\n", s, wantString)
	}

	for _, item := range items {
		if !m.Has(item.Key) {
			t.Fatalf("Map key %s not found", item.Key)
		}

		if v := m.Value(item.Key); v != item.Value {
			t.Fatalf("Wrong map item value: got %v, want %v", v, item.Value)
		}

		if i := m.indexOf(item.Key); i != item.Index {
			t.Fatalf("Wrong map item index: got %d, want %d", i, item.Index)
		}
	}

	if m.Has("invalid-key") {
		t.Fatalf("Map key %s found, want not found.", "invalid-key")
	}

	if m.Value("invalid-key") != nil {
		t.Fatalf("Map value for key %s not nil, want nil.", "invalid-key")
	}

	i := 0
	m.Iterate(func(key string, value interface{}) bool {
		item := items[i]

		if key != item.Key || value != item.Value {
			t.Fatalf("Wrong iterator item: got %s: %v, want %s: %v", key, value, item.Key, item.Value)
		}

		i++
		return true
	})

	if i != len(items) {
		t.Fatalf("Iterator not called for each map item: got %d calls, want %d", i, len(items))
	}

	replace := []struct {
		Key   string
		Value string
		Index int
	}{
		{Key: "1", Value: "bb", Index: 1},
		{Key: "3", Value: "dd", Index: 3},
	}

	for _, rep := range replace {
		m.Append(rep.Key, rep.Value)

		if !m.Has(rep.Key) {
			t.Fatalf("Map replaced item key %s not found", rep.Key)
		}

		if v := m.Value(rep.Key); v != rep.Value {
			t.Fatalf("Wrong map replaced item value: got %v, want %v", v, rep.Value)
		}

		if j := m.indexOf(rep.Key); j != rep.Index {
			t.Fatalf("Wrong map replaced item index: got %d, want %d", j, rep.Index)
		}
	}
}
