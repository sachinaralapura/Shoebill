package object

import "fmt"

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ    = "NULL"
)

type ObjecType string

// Object interface
type Object interface {
	Type() ObjecType
	Inspect() string
}

// Integer Type
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjecType { return INTEGER_OBJ }

// Boolean Type
type Boolean struct {
	Value bool
}

func (i *Boolean) Inspect() string { return fmt.Sprintf("%t", i.Value) }
func (i *Boolean) Type() ObjecType { return BOOLEAN_OBJ }

// Null type
type Null struct{}

func (n *Null) Inspect() string { return "null" }
func (n *Null) Type() ObjecType { return NULL_OBJ }
