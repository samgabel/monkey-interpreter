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
// In cases where the character is alphanumeric, the token literal can be an
// identifier/keyword or number (multi-character token literal)
//
// Before returning the token we advance our pointers into the token we will
// also advance our pointers by calling (Lexer).readChar().
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

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
	default:
		if isLetter(l.ch) {
			// grab the identifier/keyword
			tok.Literal = l.readIdentifier()
			// look up and set the TokenType const
			tok.Type = token.LookupIdent(tok.Literal)
			// (Lexer).readIdentifier() will progress the lexer's positional pointers, so we want to
			// early return here to avoid calling (Lexer).readChar() again
			return tok
		}
		if isDigit(l.ch) {
			// grabs the number
			tok.Literal = l.readNumber()
			// sets the number to an int type
			tok.Type = token.INT
			// also early return because we have numbers with multiple digits and don't want to call
			// (*Lexer).readChar() again.
			return tok
		}
		tok = newToken(token.ILLEGAL, l.ch)
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

// Will return a string of the input of an identifier or keyword (anything with contiguous letters), and will
// advance the index of the position and readPosition pointers until it encounters a non-letter character.
func (l *Lexer) readIdentifier() string {
	startPosition := l.position
	// the for loop breaks when a word ends (fn, x, ten, ...)
	for isLetter(l.ch) {
		l.readChar()
	}
	// returns a portion of the string from the start position to the position of the lexer when the for loop breaks
	return l.input[startPosition:l.position]
}

// Will return a string of the input of a number, and will advance the index of the position and readPosition pointers
// until it encounters a non-numeric character.
func (l *Lexer) readNumber() string {
	startPosition := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPosition:l.position]
}

// This helper function returns a boolean value dependent on if the ASCII code of the input byte is in the
// range of any alphabetical ASCII including the '_' character.
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// This helper function returns a boolean value dependent on if the ASCII code of the input byte is in the
// range of any numerical ASCII.
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// A helper method found in a lot of parsers used to "consume" or skip over whitespace characters.
func (l *Lexer) skipWhitespace() {
	// for any of these whitespace type characters we want to skip over them
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
