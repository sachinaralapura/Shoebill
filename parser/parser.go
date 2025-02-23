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
	OR          // ||
	AND         // &&
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX
)

// precedence table
var precedences = map[token.TokenType]int{
	token.EQUAL:     EQUALS,
	token.OR:        OR,
	token.AND:       AND,
	token.NOT_EQUAL: EQUALS,
	token.LT:        LESSGREATER,
	token.GT:        LESSGREATER,
	token.PLUS:      SUM,
	token.MINUS:     SUM,
	token.SLASH:     PRODUCT,
	token.ASTERISK:  PRODUCT,
	token.MODULO:    PRODUCT,
	token.LPAREN:    CALL,
	token.LBRACKET:  INDEX,
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

	// sendChan chan<- []byte

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
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}
	p.NextToken()
	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return stmt
}

// parse Return Statement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	returnStmt := &ast.ReturnStatement{Token: p.curToken}
	p.NextToken()
	returnStmt.ReturnValue = p.parseExpression(LOWEST)
	if p.peekTokenIs(token.SEMICOLON) {
		p.NextToken()
	}
	return returnStmt
}

// parse expression Statement
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	expression := p.parseExpression(LOWEST)
	stmt.Expression = expression
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

// parse list of comma separated expression
func (p *Parser) parseExpressionList(endToken token.TokenType) []ast.Expression {
	expressionList := []ast.Expression{}
	if p.peekTokenIs(endToken) {
		p.NextToken()
		return expressionList
	}
	p.NextToken()
	expressionList = append(expressionList, p.parseExpression(LOWEST))
	for p.peekTokenIs(token.COMMA) {
		p.NextToken()
		p.NextToken()
		expressionList = append(expressionList, p.parseExpression(LOWEST))
	}
	if !p.expectPeek(endToken) {
		return nil
	}
	return expressionList
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

// parse Identifier Expressions
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parse Integer Literal
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

// parse String Literal
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parse boolean Literal
func (p *Parser) parseBoolean() ast.Expression {
	return &ast.BoolenExpression{Token: p.curToken, Value: p.curTokenIs(token.TRUE)}
}

// parse Array Literal
func (p *Parser) parseArrayLiteral() ast.Expression {
	arrayExpression := &ast.ArrayExpression{Token: p.curToken}
	arrayExpression.Elements = p.parseExpressionList(token.RBRACKET)
	return arrayExpression
}

// parse Array Index Expression
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}
	p.NextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}
	return exp
}

// parse Hash Expression
func (p *Parser) parseHashLiteral() ast.Expression {
	hash := &ast.HashLiteral{Token: p.curToken}
	hash.Pairs = make(map[ast.Expression]ast.Expression)
	for !p.peekTokenIs(token.RBRACE) {
		p.NextToken()
		key := p.parseExpression(LOWEST)
		if !p.expectPeek(token.COLON) {
			return nil
		}
		p.NextToken()
		value := p.parseExpression(LOWEST)
		hash.Pairs[key] = value
		if !p.peekTokenIs(token.RBRACE) && !p.expectPeek(token.COMMA) {
			return nil
		}
	}
	if !p.expectPeek(token.RBRACE) {
		return nil
	}
	return hash
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

func (p *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	p.NextToken()
	expression.Condition = p.parseExpression(LOWEST)
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	expression.Consequence = p.parseBlockExpression()

	if p.peekTokenIs(token.ELSE) {
		p.NextToken()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}
		expression.Alternative = p.parseBlockExpression()
	}
	return expression
}

func (p *Parser) parseBlockExpression() *ast.BlockStatement {
	blockedExp := &ast.BlockStatement{Token: p.curToken}
	blockedExp.Statements = []ast.Statement{}
	p.NextToken()
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		statement := p.parseStatement()
		if statement != nil {
			blockedExp.Statements = append(blockedExp.Statements, statement)
		}
		p.NextToken()
	}
	return blockedExp
}

// parse Function expression
func (p *Parser) parseFunctionExpression() ast.Expression {
	fnExpression := &ast.FunctionLiteral{Token: p.curToken}
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	fnExpression.Parameters = p.parseFunctionParameters()
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	fnExpression.Body = p.parseBlockExpression()
	return fnExpression
}

func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}
	if p.peekTokenIs(token.RPAREN) {
		p.NextToken()
		return identifiers
	}
	p.NextToken()
	identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, identifier)
	for p.peekTokenIs(token.COMMA) {
		p.NextToken()
		p.NextToken()
		identifier := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, identifier)
	}
	if !p.expectPeek(token.RPAREN) {
		return nil
	}
	return identifiers
}

// parse Call Expressions
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	callExp := &ast.CallExpression{Token: p.curToken, Function: function}
	callExp.Arguments = p.parseExpressionList(token.RPAREN)
	return callExp
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
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionExpression)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseHashLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.MODULO, p.parseInfixExpression)
	p.registerInfix(token.EQUAL, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQUAL, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)

	// read next two token so the curToken and peekToken are read
	p.NextToken()
	p.NextToken()

	return p
}
