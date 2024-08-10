package ast

import "github.com/samgabel/monkey-interpreter/token"

// Every node in our AST has to implement the Node interface, meaning that they have to provide
// a TokenLiteral() method that returns the literal value of the token.
//
// TokenLiteral() will only be used for debugging and testing.
type Node interface {
	TokenLiteral() string
}

// Statements do NOT produce values, like (let x = 5)
//
// statementNode() serves as a dummy method to help us guide the Go compiler and throw errors
type Statement interface {
	statementNode()
	Node // aka needs TokenLiteral() method as well to implement this interface
}

// Expressions produce values, like (return 5)
//
// expressionNode() serves as a dummy method to help us guide the Go compiler and throw errors
type Expression interface {
	expressionNode()
	Node // aka needs TokenLiteral() method as well to implement this interface
}

// The Program struct implements the Node interface and is our root node of every AST.
//
// Every Monkey program is simply just a series of statements contained in the Program.Statements,
// which is just a slice of AST nodes that implement the Statement interface.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

// Implements the Statement interface (which implements the Node interface), this is a struct that
// contains the token.LET token, the associated identifier and it's expression.
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

// for implementation puposes only
func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

// Implements the Expression interface (which implements the Node interface), this is a struct that
// containes the token.IDENT token, and its associated value string.
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

// for implementation puposes only
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
