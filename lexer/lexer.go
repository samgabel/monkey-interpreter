package lexer

import "github.com/samgabel/monkey-interpreter/token"

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

// This function servers to construct a new instance of a Lexer and initialize
// our position in the input.
//
// By calling (Lexer).readChar() we will initialize our Lexer ch, position, and
// readPosition
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// This method will look under the current character under
// examination (l.ch) and return a token depending on which character it is.
//
// Before returning the token we advance our pointers into the token we will
// also advance our pointers by calling (Lexer).readChar().
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	}

	l.readChar()
	return tok
}

// This is a helper function in order to format the token.Token type to be returned
func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

// This method's purpose is to give us the next character and to advance our position in the input string.
func (l *Lexer) readChar() {
	// check if we have reached the end of the input
	if l.readPosition >= len(l.input) {
		l.ch = 0 // 0 is the ASCII code for the "NUL" character
	} else {
		l.ch = l.input[l.readPosition]
	}

	// set our new current position to be the previous readPosition and increment the new readPosition by 1.
	l.position = l.readPosition
	l.readPosition += 1
}
