package parser

import (
	"fmt"

	"github.com/sachinaralapura/shoebill/ast"
	"github.com/sachinaralapura/shoebill/lexer"
	"github.com/sachinaralapura/shoebill/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression // left expression as argument
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	sendChan chan<- []byte

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// helper functions

func (p *Parser) registerPrefix(tokenType token.TokenType, prefixFn prefixParseFn) {
	p.prefixParseFns[tokenType] = prefixFn
}

func (p *Parser) registerInfix(tokenType token.TokenType, infixFn infixParseFn) {
	p.infixParseFns[tokenType] = infixFn
}

// parse programm
func (p *Parser) ParseProgram() *ast.Program {
	// create a program ( root node )
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// append statement to the program statement
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.NextToken()
	}
	return program
}

// ----- parse Statemets -----
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// ----------------- parse Let Statements -----------------
func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	// check if next token is Identifier
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// TODO : Evaluate the expression
	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return stmt
}

// ----------------- parse Return Statement -----------------
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStmt := &ast.ReturnStatement{Token: p.curToken}

	// TODO : Evaluate the expression
	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return returnStmt
}

// ----------------- parse expression Statement -----------------
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expression := p.parseExpression(LOWEST)
	stmt := &ast.ExpressionStatement{Token: p.curToken, Expression: expression}
	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return stmt
}

// ----------------- parse Expression -----------------
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// -------------------------------------------------------
func (p *Parser) NextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.NextToken()
		return true
	} else {
		p.peekErrors(t)
		return false
	}
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekErrors(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s , got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) Errors() []string {
	return p.errors
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}
	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)

	// read next two token so the curToken and peekToken are read
	p.NextToken()
	p.NextToken()

	return p
}
