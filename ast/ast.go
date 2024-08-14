package ast

import (
	"strings"

	"github.com/samgabel/monkey-interpreter/token"
)

// Every node in our AST has to implement the Node interface, meaning that they have to provide
// a TokenLiteral() method that returns the literal value of the token.
//
// TokenLiteral() and String() will only be used for debugging and testing.
type Node interface {
	TokenLiteral() string
	String() string
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

// Returns the value of each statement's String() method and aggregates them into a buffer
// that is then returned. Great for testing and debugging the contents of the Program.
func (p *Program) String() string {
	var out strings.Builder

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
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

// This simply will construct a string of the entire statement: "let <identifier> = <expression>;"
func (ls *LetStatement) String() string {
	var out strings.Builder

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

// Implements the Expression interface (which implements the Node interface), this is a struct that
// containes the token.IDENT token, and its associated value string.
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

// for implementation puposes only
func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// Implements the Statement interface (which implements the Node interface), this is a struct that
// contains the token.RETURN token, and the associated expression.
type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

// for implementation purposes only
func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// This simply will construct a string of the entire statement: "return <expression>;"
func (rs *ReturnStatement) String() string {
	var out strings.Builder

	out.WriteString(rs.TokenLiteral() + " ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

// Implements the Statement interface (which implements the Node interface). We have "expression statements" in
// this language because we wan't to allow the use of standalone Expressions (or Unused Expressions).
type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

// for implementation purposes only
func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }

// This simply will construct a string of the entire statement: "<expression>"
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}
