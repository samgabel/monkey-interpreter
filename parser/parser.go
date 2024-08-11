package parser

import (
	"fmt"

	"github.com/samgabel/monkey-interpreter/ast"
	"github.com/samgabel/monkey-interpreter/lexer"
	"github.com/samgabel/monkey-interpreter/token"
)

// We will repeatedly call (*Lexer).NextToken() in order to get tokens from the input.
//
// curToken and peekToken act exactly like the two "pointers" our lexer has (position nad readPosition).
// Intstead of pointing to a character in the input they point to the current and the next token.
//
// We need the curToken in order to decide what to do next and the peekToken for this decision if curToken
// doesn't give us enough information.
type Parser struct {
	lexer *lexer.Lexer

	errors []string

	curToken  token.Token
	peekToken token.Token
}

// This function servers to construct a new instance of a Parser.
//
// By calling (*Parser).nextToken() we will initialize our curToken and peekToken.
func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer: l,
		errors: []string{},
	}

	// read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

// This Parser method will return a slice of all the accumulated errors so far.
func (p *Parser) Errors() []string{
	return p.errors
}

// This Parser method will be used to add an error to (Parser).errors when the type of peekToken
// doesn't match the expectation.
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// This Parser method will set the curToken to the previous peekToken and then call the (*Lexer).nextToken()
// method in order to "read" forward in the lexer token queue.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
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
		return nil
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
