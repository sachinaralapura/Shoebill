/*
Tokens are the output of the lexical analyzer and serve as the input to the parser
Each token is classified into a type, such as:

	Keywords: if, while, for
	Identifiers: variable names like x, counter
	Literals: 42, 3.14, 'c'
	Operators: +, -, *, =
	Delimiters: (, ), {, }
*/
package token

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//Identified + literals
	IDENT  = "IDENT"
	INT    = "INT"
	STRING = "STRING"

	//Operators
	ASSIGN    = "="
	PLUS      = "+"
	MINUS     = "-"
	ASTERISK  = "*"
	SLASH     = "/"
	MODULO    = "%"
	BANG      = "!"
	LT        = "<"
	GT        = ">"
	EQUAL     = "=="
	NOT_EQUAL = "!="

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"
	QUOTE     = "\""

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	CONST    = "CONST"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	IF       = "IF"
	ELSE     = "ELSE"
	FOR      = "FOR"
	WHILE    = "WHILE"
	RETURN   = "RETURN"
	CLASS    = "CLASS"
)

var Keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"const":  CONST,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"class":  CLASS,
}

func LookUpIdent(ident string) TokenType {
	if tok, ok := Keywords[ident]; ok {
		return tok
	}
	return IDENT
}

type TokenType string

type Token struct {
	Type    TokenType // token name
	Literal string    // token attribute
	Line    int       // line Number
}
