package object

import (
	"bytes"
	"fmt"
	"strings"
)

type HashKey struct {
	Type  ObjecType
	Value uint64
}

type Hashable interface {
	HashKey() HashKey
}

type HashPair struct {
	Key   Object
	Value Object
}

// Hash Type object
// Implements object and Hashable interface
type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjecType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s:%s", pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
