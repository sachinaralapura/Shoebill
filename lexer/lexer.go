/*
A lexeme is the actual sequence of characters in the source code that matches a pattern for a token.
Lexemes are the input to the lexer
*/
package lexer

import (
	"bytes"
	"fmt"
	"regexp"

	"github.com/sachinaralapura/shoebill/token"
)

type Buffer struct {
	buffer1 []rune
	buffer2 []rune
	CurBuf  *[]rune
}

type Lexer struct {
	// double buffer
	Buffer

	position     int
	readPosition int
	ch           rune

	recieveChan <-chan []byte

	currentLineNumber int
	tokens            []token.Token
}

func (l *Lexer) String() string {
	var out bytes.Buffer
	for _, tok := range l.tokens {
		out.WriteString(fmt.Sprintf("<%s, %s , [%d]>\n", tok.Type, tok.Literal, tok.Line))
	}
	return out.String()
}

// LoadBuffer reads a chunk of data from the recieveChan
func (l *Lexer) LoadBuffer() bool {

	// recieve chunk from the channel
	data, ok := <-l.recieveChan
	if !ok {
		return false
	}
	runes := []rune(string(data))

	// swap buffers
	if l.CurBuf == &l.buffer1 {
		l.buffer2 = runes
		l.CurBuf = &l.buffer2
	} else {
		l.buffer1 = runes
		l.CurBuf = &l.buffer1
	}

	l.position = 0
	l.readPosition = 0

	// read first character
	l.ch = (*l.CurBuf)[l.readPosition]
	l.position = l.readPosition
	l.readPosition++
	return true
}

// readChar reads the next character from the buffer1 and advances the read position.
// If the end of the buffer1 is reached, it sets the current character to 0.
func (l *Lexer) readChar() {
	if l.readPosition >= len(*l.CurBuf) {
		if l.LoadBuffer() {
			l.ch = (*l.CurBuf)[l.readPosition]
		} else {
			l.ch = 0
		}
	} else {
		l.ch = (*l.CurBuf)[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// peekChar returns the next character in the buffer1 without advancing the read position.
// It returns 0 if the end of the buffer1 is reached.
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(*l.CurBuf) {
		return 0
	}
	return (*l.CurBuf)[l.readPosition]
}

// readIdentifier reads an identifier from the buffer1
func (l *Lexer) readIdentifier() []rune {
	var currentBuffer *[]rune = l.CurBuf
	initPos := l.position
	for isIdentifierCharacter(l.ch) {
		l.readChar()
	}
	if currentBuffer != l.CurBuf {
		prevBufPart := (*currentBuffer)[initPos:]
		currBufPart := (*l.CurBuf)[0:l.position]
		return append(prevBufPart, currBufPart...)
	}

	return (*l.CurBuf)[initPos:l.position]
}

// readNumber reads a number from the buffer1
//
// It reads characters until it encounters a non-digit.
func (l *Lexer) readNumber() []rune {
	var currentBuffer *[]rune = l.CurBuf
	initPos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	if currentBuffer != l.CurBuf {
		prevBufPart := (*currentBuffer)[initPos:]
		currBufPart := (*l.CurBuf)[0:l.position]
		return append(prevBufPart, currBufPart...)
	}
	return (*l.CurBuf)[initPos:l.position]
}

func (l *Lexer) readString() []rune {
	l.readChar() // skip '"'
	var currentBuffer *[]rune = l.CurBuf
	initPos := l.position
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	if currentBuffer != l.CurBuf {
		prevBufPart := (*currentBuffer)[initPos:]
		currBufPart := (*l.CurBuf)[0:l.position]
		return append(prevBufPart, currBufPart...)
	}
	return (*l.CurBuf)[initPos:l.position]
}

// skipWhiteSpace skips all whitespace characters in the input.
func (l *Lexer) skipWhiteSpace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	if l.ch == '\n' {
		l.currentLineNumber += 1
	}
	l.skipWhiteSpace()
	curChar := string(l.ch)
	switch l.ch {
	case '=':
		if l.peekChar() == l.ch {
			l.readChar()
			tok = newToken(token.EQUAL, curChar+string(l.ch), l.currentLineNumber)
		} else {
			tok = newToken(token.ASSIGN, curChar, l.currentLineNumber)
		}
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = newToken(token.NOT_EQUAL, curChar+string(l.ch), l.currentLineNumber)
		} else {
			tok = newToken(token.BANG, curChar, l.currentLineNumber)
		}
	case '&':
		if l.peekChar() == '&' {
			l.readChar()
			tok = newToken(token.AND, "&&", l.currentLineNumber)
		}
	case '|':
		if l.peekChar() == '|' {
			l.readChar()
			tok = newToken(token.OR, "||", l.currentLineNumber)
		}
	case '"':
		tok.Type = token.STRING
		tok.Line = l.currentLineNumber
		tok.Literal = string(l.readString())
	case ';':
		tok = newToken(token.SEMICOLON, curChar, l.currentLineNumber)
	case '(':
		tok = newToken(token.LPAREN, curChar, l.currentLineNumber)
	case ')':
		tok = newToken(token.RPAREN, curChar, l.currentLineNumber)
	case ',':
		tok = newToken(token.COMMA, curChar, l.currentLineNumber)
	case '{':
		tok = newToken(token.LBRACE, curChar, l.currentLineNumber)
	case '}':
		tok = newToken(token.RBRACE, curChar, l.currentLineNumber)
	case '[':
		tok = newToken(token.LBRACKET, curChar, l.currentLineNumber)
	case ']':
		tok = newToken(token.RBRACKET, curChar, l.currentLineNumber)
	case '+':
		tok = newToken(token.PLUS, curChar, l.currentLineNumber)
	case '-':
		tok = newToken(token.MINUS, curChar, l.currentLineNumber)
	case '*':
		tok = newToken(token.ASTERISK, curChar, l.currentLineNumber)
	case '/':
		tok = newToken(token.SLASH, curChar, l.currentLineNumber)
	case '%':
		tok = newToken(token.MODULO, curChar, l.currentLineNumber)
	case '<':
		tok = newToken(token.LT, curChar, l.currentLineNumber)
	case '>':
		tok = newToken(token.GT, curChar, l.currentLineNumber)

	case 0:
		tok = newToken(token.EOF, "", l.currentLineNumber)
	default:
		if isLetter(l.ch) {
			// check if it's keyword , identifier
			identifier := string(l.readIdentifier())
			tokenType := token.LookUpIdent(identifier)
			token := newToken(tokenType, identifier, l.currentLineNumber)

			l.addToken(token) // add token to l.tokens
			return token
		} else if isDigit(l.ch) {
			tok.Literal = string(l.readNumber())
			tok.Type = token.INT
			tok.Line = l.currentLineNumber
			l.addToken(tok)
			return tok
		} else {
			tok = newToken(token.ILLEGAL, curChar, l.currentLineNumber)
		}
	}

	// read next char
	l.readChar()
	l.addToken(tok)
	return tok
}

func (l *Lexer) addToken(token token.Token) {
	l.tokens = append(l.tokens, token)
}

// ------------------ Helper functions -----------------------------

func New(recieve <-chan []byte) *Lexer {
	l := &Lexer{recieveChan: recieve, currentLineNumber: 1}
	l.CurBuf = &l.buffer2
	return l
}

func NewFromString(input string) *Lexer {
	l := &Lexer{}
	l.buffer1 = []rune(input)
	l.readChar()
	return l
}

func isLetter(char rune) bool {
	return regexp.MustCompile("^[a-zA-z_]$").MatchString(string(char))
}

func isIdentifierCharacter(char rune) bool {
	return regexp.MustCompile("^[a-zA-Z_0-9]$").MatchString(string(char))
}

// isAcceptableIdentifierCharacter checks if the character is acceptable within a broader context of identifier if you wish to extend it later.
// func isAcceptableStringCharacter(char rune) bool {
// 	return regexp.MustCompile("^[^\"]$").MatchString(string(char))
// }

func isDigit(char rune) bool {
	return regexp.MustCompile("^[0-9.]$").MatchString(string(char))
}

func newToken(tokenType token.TokenType, attribute string, lineNumber int) token.Token {
	return token.Token{Type: tokenType, Literal: attribute, Line: lineNumber}
}

// func isIdentifier(input []rune) bool {
// 	return regexp.MustCompile("^[a-zA-Z_][a-zA-Z0-9_]{0,19}$").MatchString(string(input))
// }
