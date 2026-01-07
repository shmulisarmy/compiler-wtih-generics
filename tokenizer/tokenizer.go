package tokenizer

import (
	"fmt"
	"unicode"
)

type TokenType string

const (
	// Literals
	IDENT  TokenType = "IDENT"
	INT    TokenType = "INT"
	STRING TokenType = "STRING"

	// Keywords
	FUNC   TokenType = "FUNC"
	CLASS  TokenType = "CLASS"
	IF     TokenType = "IF"
	WHILE  TokenType = "WHILE"
	RETURN TokenType = "RETURN"
	ELSE   TokenType = "ELSE"
	USING  TokenType = "USING"
	BREAK  TokenType = "BREAK"
	UNLESS TokenType = "UNLESS"

	// Operators
	ASSIGN TokenType = "ASSIGN" // =
	PLUS   TokenType = "PLUS"   // +
	MINUS  TokenType = "MINUS"  // -
	MUL    TokenType = "MUL"    // *
	DIV    TokenType = "DIV"    // /
	EQ     TokenType = "EQ"     // ==
	NEQ    TokenType = "NEQ"    // !=
	GT     TokenType = "GT"     // >
	LT     TokenType = "LT"     // <
	GTE    TokenType = "GTE"    // >=
	LTE    TokenType = "LTE"    // <=
	AND    TokenType = "AND"    // &&
	OR     TokenType = "OR"     // ||
	NOT    TokenType = "NOT"    // !

	//  Delimiters
	LPAREN TokenType = "LPAREN" // (
	RPAREN TokenType = "RPAREN" // )
	LBRACE TokenType = "LBRACE" // {
	RBRACE TokenType = "RBRACE" // }
	DOT    TokenType = "DOT"    // .
	COMMA  TokenType = "COMMA"  // ,
	COLON  TokenType = "COLON"  // :

	// Special
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"
)

var keywords = map[string]TokenType{
	"func":   FUNC,
	"class":  CLASS,
	"if":     IF,
	"while":  WHILE,
	"return": RETURN,
	"else":   ELSE,
	"break":  BREAK,
	"unless": UNLESS,
	"using":  USING,
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

type Tokenizer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

func New(input string) *Tokenizer {
	t := &Tokenizer{
		input:  input,
		line:   1,
		column: 0,
	}
	t.readChar()
	return t
}

func (t *Tokenizer) readChar() {
	if t.in_range() {
		t.ch = 0
	} else {
		t.ch = t.input[t.readPosition]
	}
	t.position = t.readPosition
	t.readPosition++
	t.column++

	if t.ch == '\n' {
		t.line++
		t.column = 0
	}
}

func (t *Tokenizer) in_range() bool {
	return t.readPosition >= len(t.input)
}

func (t *Tokenizer) peekChar() byte {
	if t.in_range() {
		return 0
	}
	return t.input[t.readPosition]
}

func (t *Tokenizer) NextToken() Token {
	var tok Token

	t.skipWhitespace()

	tok.Line = t.line
	tok.Column = t.column

	switch t.ch {
	case '=':
		if t.peekChar() == '=' {
			t.readChar()
			tok = t.newToken(EQ, string(t.ch))
		} else {
			tok = t.newToken(ASSIGN, string(t.ch))
		}
	case '(':
		tok = t.newToken(LPAREN, string(t.ch))
	case ')':
		tok = t.newToken(RPAREN, string(t.ch))
	case '{':
		tok = t.newToken(LBRACE, string(t.ch))
	case '}':
		tok = t.newToken(RBRACE, string(t.ch))
	case '.':
		tok = t.newToken(DOT, string(t.ch))
	case ',':
		tok = t.newToken(COMMA, string(t.ch))
	case ':':
		tok = t.newToken(COLON, string(t.ch))
	//operators
	case '+':
		tok = t.newToken(PLUS, string(t.ch))
	case '-':
		tok = t.newToken(MINUS, string(t.ch))
	case '*':
		tok = t.newToken(MUL, string(t.ch))
	case '/':
		tok = t.newToken(DIV, string(t.ch))

	case '!':
		tok = t.newToken(NOT, string(t.ch))
	case '<':
		tok = t.newToken(LT, string(t.ch))
	case '>':
		tok = t.newToken(GT, string(t.ch))
	case '&':
		tok = t.newToken(AND, string(t.ch))
	case '|':
		tok = t.newToken(OR, string(t.ch))
	case '"':
		tok.Type = STRING
		tok.Literal = t.readString()

		return tok
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(t.ch) {
			tok.Literal = t.readIdentifier()
			tok.Type = lookupIdent(tok.Literal)
			return tok
		} else if isDigit(t.ch) {
			tok.Type = INT
			tok.Literal = t.readNumber()
			return tok
		} else {
			tok = t.newToken(ILLEGAL, string(t.ch))
		}
	}

	t.readChar()
	return tok
}

func (t *Tokenizer) Tokenize() []Token {
	var tokens []Token

	for {
		tok := t.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}

	return tokens
}

func (t *Tokenizer) newToken(tokenType TokenType, literal string) Token {
	return Token{
		Type:    tokenType,
		Literal: literal,
		Line:    t.line,
		Column:  t.column,
	}
}

func (t *Tokenizer) skipWhitespace() {
	for t.ch == ' ' || t.ch == '\t' || t.ch == '\n' || t.ch == '\r' {
		t.readChar()
	}
}

func (t *Tokenizer) readIdentifier() string {
	position := t.position
	for isLetter(t.ch) || isDigit(t.ch) {
		t.readChar()
	}
	return t.input[position:t.position]
}

func (t *Tokenizer) readNumber() string {
	position := t.position
	for isDigit(t.ch) {
		t.readChar()
	}
	return t.input[position:t.position]
}

func (t *Tokenizer) readString() string {
	position := t.position + 1
	for {
		t.readChar()
		if t.ch == '"' || t.ch == 0 {
			break
		}
	}
	token := t.input[position:t.position]
	t.readChar()
	return token
}

func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func lookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %s, Literal: %q, Line: %d, Column: %d}", t.Type, t.Literal, t.Line, t.Column)
}
