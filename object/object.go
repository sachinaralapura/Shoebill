package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/sachinaralapura/shoebill/ast"
)

const (
	INTEGER_OBJ  = "INTEGER"
	STRING_OBJ   = "STRING"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN"
	ERROR_OBJ    = "ERROR"
	FUCNTION_OBJ = "FUNCTION"
	BUILDIN_OBJ  = "BUILDIN"
	ARRAY_OBJ    = "ARRAY"
)

type ObjecType string
type BuildInFunc func(args ...Object) Object

// Object interface
type Object interface {
	Type() ObjecType
	Inspect() string
}

// build in function Object
type BuildIn struct {
	Value BuildInFunc
}

func (bi *BuildIn) Inspect() string { return BUILDIN_OBJ + "function" }
func (bi *BuildIn) Type() ObjecType { return BUILDIN_OBJ }

// Integer Type Object
type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjecType { return INTEGER_OBJ }

// String Type Object
type String struct {
	Value string
}

func (s *String) Inspect() string { return s.Value }
func (s *String) Type() ObjecType { return STRING_OBJ }

// Boolean Type Object
type Boolean struct {
	Value bool
}

func (i *Boolean) Inspect() string { return fmt.Sprintf("%t", i.Value) }
func (i *Boolean) Type() ObjecType { return BOOLEAN_OBJ }

// Array Type Object
type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjecType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ","))
	out.WriteString("]")
	return out.String()
}

// Null Type Object
type Null struct{}

func (n *Null) Inspect() string { return "null" }
func (n *Null) Type() ObjecType { return NULL_OBJ }

// Return Object
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
