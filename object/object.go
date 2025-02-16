package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sachinaralapura/shoebill/ast"
)

const (
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN"
	ERROR_OBJ    = "ERROR"
	FUCNTION_OBJ = "FUNCTION"
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

// Return type
type Return struct {
	Value Object
}

func (r *Return) Inspect() string { return r.Value.Inspect() }
func (r *Return) Type() ObjecType { return RETURN_OBJ }

// Function Literal Object
type FunctionObject struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *FunctionObject) Type() ObjecType { return FUCNTION_OBJ }
func (f *FunctionObject) Inspect() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}
	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")
	return out.String()
}

// Error Object
type ErrorObject struct {
	Message string
}

func (e *ErrorObject) Inspect() string { return e.Message }
func (e *ErrorObject) Type() ObjecType { return ERROR_OBJ }
