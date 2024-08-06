package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	ILLEGAL = "ILLEGAL" // signifies an illegal token/character we don't know about
	EOF     = "EOF"     // tells our parser when we can stop

	// Identifiers + literals
	IDENT = "IDENT" // add, foobar, x, y, ...
	INT   = "INT"   // 123456

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	// Special
	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

// Our map to store our Keywords
var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

// Returns a keyword TokenType const (FUNCTION, LET) if found in the keywords map or
// just the IDENT (user defined identifiers) TokenType const if not found in the map.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
