// Package ordered implements a simple ordered map to ensure that the insertion order
// in the elements of a map is maintained.
package ordered

import (
	"fmt"
	"strings"
)

// pair represents a kay-value pair into the ordered map.
type pair struct {
	key   string
	value interface{}
}

// mapIterator defines the iterator function signature to be used on Map.Iterate.
type mapIterator func(key string, value interface{}) bool

// Map represents an ordered map.
type Map struct {
	items []pair // Ordered map items.
}

// NewMap creates a new ordered map object.
func NewMap() *Map {
	return &Map{
		items: make([]pair, 0, 16),
	}
}

// String implements the stringer interface.
func (m *Map) String() string {
	var builder strings.Builder
	builder.WriteString("ordered.Map[")

	for _, item := range m.items {
		builder.WriteString(fmt.Sprintf("%s:%v ", item.key, item.value))
	}

	return strings.TrimRight(builder.String(), " ") + "]"
}

// Len returns the number of kay-value pairs in the map.
func (m *Map) Len() int {
	return len(m.items)
}

// Has checks if map contains given key.
func (m *Map) Has(key string) bool {
	for _, item := range m.items {
		if item.key == key {
			return true
		}
	}

	return false
}

// Value returns the value for given key, nil if key don't exist.
func (m *Map) Value(key string) interface{} {
	for _, item := range m.items {
		if item.key == key {
			return item.value
		}
	}

	return nil
}

// Append appends a new key-value pair to the map. If the key already exists the value will be overwritten.
func (m *Map) Append(key string, value interface{}) {
	if i := m.indexOf(key); i >= 0 {
		m.items[i].value = value
	} else {
		m.items = append(m.items, pair{key: key, value: value})
	}
}

// Iterate iterates through the map key-value pairs. If the passed function returns false the iteration will
// stop and Iterate returns immediately.
func (m *Map) Iterate(iterator mapIterator) {
	for _, item := range m.items {
		if !iterator(item.key, item.value) {
			break
		}
	}
}

// indexOf return the zero-based index of given key, -1 if key do not exist.
func (m *Map) indexOf(key string) int {
	for i, item := range m.items {
		if item.key == key {
			return i
		}
	}

	return -1
}
