package ast

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/sachinaralapura/shoebill/token"
)

// -------------------- interfaces ------------------
type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// ------------------------- program -----------------
// Node Interface
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		_, err := out.WriteString(s.String())
		out.WriteString("\n")
		if err != nil {
			log.Println(err)
		}
	}
	return out.String()
}

/*
--------------------- Statements ------------------------------
--------------------- Let statement --------------------
statement interface
*/
type LetStatement struct {
	Token token.Token // 'let' token
	Name  *Identifier // identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString("=")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	out.WriteString(";")
	return out.String()
}

// --------------------- Return statement --------------------
// statement interface
type ReturnStatement struct {
	Token       token.Token // 'return' token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }
func (rs *ReturnStatement) String() string {
	var out bytes.Buffer
	out.WriteString(rs.TokenLiteral() + " ")
	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}
	out.WriteString(";")
	return out.String()
}

// --------------------- Expression statement --------------------
// statement interface
type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// --------------------- Block statements --------------------
// statement interface
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{")
	for _, statement := range bs.Statements {
		out.WriteString(statement.String())
	}
	out.WriteString("}")
	return out.String()
}

// ---------------------- Expressions -------------------

/*
Identifier
implements expression interface

	It represent "Identifiers" in AST ex:
	>> let five = 5;
*/
type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

/*
Implements expression interface

	Integer Literal expression ex :
	>> 5;
	>> 4;
	>> 34;
*/
type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntegerLiteral) expressionNode()      {}
func (i *IntegerLiteral) TokenLiteral() string { return i.Token.Literal }
func (i *IntegerLiteral) String() string       { return fmt.Sprint(i.Value) }

// implements expression interface
type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Value }

/*
Boolen expresssion
Implements expression interface
*/
type BoolenExpression struct {
	Token token.Token
	Value bool
}

func (be *BoolenExpression) expressionNode()      {}
func (be *BoolenExpression) TokenLiteral() string { return be.Token.Literal }
func (be *BoolenExpression) String() string       { return be.Token.Literal }

// Array Expression
// implement Expression interface
type ArrayExpression struct {
	Token    token.Token
	Elements []Expression
}

func (ae *ArrayExpression) expressionNode()      {}
func (ae *ArrayExpression) TokenLiteral() string { return ae.Token.Literal }
func (ae *ArrayExpression) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, ele := range ae.Elements {
		elements = append(elements, ele.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ","))
	out.WriteString("]")
	return out.String()
}

// Index Expression
// <expression>[<expression>]
type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("]")
	return out.String()
}

// Hash Literal
// {<expression>:<expression> , <expression> : <expression>}
type HashLiteral struct {
	Token token.Token
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	var out bytes.Buffer
	pairs := []string{}
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}

/*
IF ELSE Expression
Implements Expression interface
*/
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (ie *IfExpression) expressionNode()      {}
func (ie *IfExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IfExpression) String() string {
	var out bytes.Buffer
	out.WriteString("if(")
	out.WriteString(ie.Condition.String())
	out.WriteString(")")
	out.WriteString(ie.Consequence.String())

	if ie.Alternative != nil {
		out.WriteString("else")
		out.WriteString(ie.Alternative.String())
	}
	return out.String()
}

/*
Function Literal
Implements Expression
*/
type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fe *FunctionLiteral) TokenLiteral() string { return fe.Token.Literal }
func (fe *FunctionLiteral) expressionNode()      {}
func (fe *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fe.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fe.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(")")
	out.WriteString(fe.Body.String())

	return out.String()
}

/*
function call expression
implements Expression interface
*/
type CallExpression struct {
	Token     token.Token
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString((strings.Join(args, ",")))
	out.WriteString(")")

	return out.String()
}

/*
Prefix Expression , Implements Expressoin interface
-5 or !true
*/
type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

/*
Infix Expression , Implements Expression interface
<Expression> <Infix Operator> <Expression>
*/
type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(ie.Operator)
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}
