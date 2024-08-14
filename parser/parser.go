package parser

import (
	"fmt"

	"github.com/samgabel/monkey-interpreter/ast"
	"github.com/samgabel/monkey-interpreter/lexer"
	"github.com/samgabel/monkey-interpreter/token"
)

// Constants are ordered by precedence (higher int = higher precedence)
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
	// Function signature that gets called when we encounter the associated token type in prefix position.
	prefixParseFn func() ast.Expression
	// Function signature that gets called when we encounter the associated token type in infix position.
	//
	// We take an additional ast.Expression as input because we have to deal with the "left side" of the
	// definition as well.
	infixParseFn func(ast.Expression) ast.Expression
)

// We will repeatedly call (*Lexer).NextToken() in order to get tokens from the input.
//
// curToken and peekToken act exactly like the two "pointers" our lexer has (position nad readPosition).
// Intstead of pointing to a character in the input they point to the current and the next token.
// We need the curToken in order to decide what to do next and the peekToken for this decision if curToken
// doesn't give us enough information.
//
// prefix/infixParseFns maps serve to check if our curToken.Type has an associated prefix/infix parsing function.
type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

// This function servers to construct a new instance of a Parser.
//
// By calling (*Parser).nextToken() we will initialize our curToken and peekToken.
func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:  l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)

	// read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// This Parser method will set the curToken to the previous peekToken and then call the (*Lexer).nextToken()
// method in order to "read" forward in the lexer token queue.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

// Boolean value to match our curToken to our input token Type.
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// Boolean value to match our peekToken to our input token Type.
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Will advance our Parser tokens if our peekToken matches our input token Type.
//
// This is considered an "assertion function" and nearly all parsers share this.
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// This Parser method will return a slice of all the accumulated errors so far.
func (p *Parser) Errors() []string {
	return p.errors
}

// This Parser method will be used to add an error to (Parser).errors when the type of peekToken
// doesn't match the expectation.
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// This will be the entry point of the Recursive-Descent Parser.
//
// It will construct the root node of the AST and then iterate over every token by calling (*Parser).nextToken() repeatedly.
// This advances the Parser's curToken and peekToken.
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}

	// iterate over all tokens till we hit the EOF token
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

// This is an intermediate method that handles parsing of statements based on the curToken Type.
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

// This method handles statement processing specifically for LET statements.
// "let <identifier> = <expression>;"
func (p *Parser) parseLetStatement() *ast.LetStatement {
	// our curToken at this point should be a LET token
	stmt := &ast.LetStatement{Token: p.curToken}

	// if our next token is not an identifier, we should not return our LetStatement because it is flawed
	// advances our Parser tokens if it passes
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	// populate our statement identifier with the curToken, and assign it's pointer to our statement
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// same as above but after our identifier, our curToken should be an assignment operator
	// advances our Parser tokens if it passes
	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	// TODO:: Skip expressions for now:
	// we're skipping the expressions until we encounter a semicolon.
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// This method handles statement processing specifically for RETURN statements.
// "return <expression>;"
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	// our curToken at this point should be a RETURN token
	stmt := &ast.ReturnStatement{Token: p.curToken}

	// go to the next token which should be our Expression
	p.nextToken()

	// TODO:: Skip expressions for now:
	// we're skipping the expressions until we encounter a semicolon.
	for !p.curTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// This method handles statement processing specifically for standalone Expressions (aka Expression Statements)
// "<expression>"
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	// check and handle OPTIONAL semicolon (makes it easy to type in things like "5 + 5" directly into the REPL)
	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// TODO: Flesh out
// our first rendition will only check whether we have a parsing function associated with p.curToken.Type in the
// prefix postition; this is enough to pass TestIdentifierExpression().
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		return nil
	}
	leftExp := prefix()

	return leftExp
}

// This prefix parse function handles parsing of Identifiers. It happens to be that this is very simple,
// we only need to return a pointer to the curToken wrapped as an ast.Identifier.
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// A helper function that adds entries to the (*Parser).prefixParseFns map
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// A helper function that adds entries to the (*Parser).infixParseFns map
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
