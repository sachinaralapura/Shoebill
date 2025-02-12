package parser

import (
	"fmt"
	"strconv"

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

// precedence table
var precedences = map[token.TokenType]int{
	token.EQUAL:     EQUALS,
	token.NOT_EQUAL: EQUALS,
	token.LT:        LESSGREATER,
	token.GT:        LESSGREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.SLASH:     PRODUCT,
	token.ASTERISK:  PRODUCT,
}

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

// determine the precedence of next token
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// determine the precedence of the current token
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
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

// parse Let Statements
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

// parse Return Statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStmt := &ast.ReturnStatement{Token: p.curToken}

	// TODO : Evaluate the expression
	for !p.curTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return returnStmt
}

// parse expression Statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	expression := p.parseExpression(LOWEST)
	stmt := &ast.ExpressionStatement{Token: p.curToken, Expression: expression}
	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return stmt
}

// ---------------- parse Expression --------------------
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}

	leftExp := prefix()
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}
		p.NextToken()
		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// parse Identifier Expressions
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parse Integer Expression
func (p *Parser) parseIntegerExpression() ast.Expression {
	integerLiteral := &ast.IntegerLiteral{Token: p.curToken}
	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	integerLiteral.Value = value
	return integerLiteral
}

// parse boolean expression
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BoolenExpression{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// parse grouped Expression
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.NextToken()
	exp := p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return exp
}

// parse Prefix Expression
func (p *Parser) parsePrefixExpression() ast.Expression {
	prefixExp := &ast.PrefixExpression{Token: p.curToken, Operator: p.curToken.Literal}
	p.NextToken()
	exp := p.parseExpression(PREFIX)
	prefixExp.Right = exp
	return prefixExp
}

// parse Infix Expression
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infixExp := &ast.InfixExpression{Token: p.curToken, Operator: p.curToken.Literal, Left: left}
	precedence := p.curPrecedence()
	p.NextToken()
	infixExp.Right = p.parseExpression(precedence)
	return infixExp
}

// proceeds to next token
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
	p.registerPrefix(token.INT, p.parseIntegerExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.EQUAL, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)

	// read next two token so the curToken and peekToken are read
	p.NextToken()
	p.NextToken()

	return p
}
