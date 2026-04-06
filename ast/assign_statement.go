package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

type AssignStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (s *AssignStatement) statementNode() {}

// TokenLiteral returns the literal token.
func (s *AssignStatement) TokenLiteral() string {
	return s.Token.Literal
}

// String returns this object as a string.
func (as *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Name.String())
	out.WriteString(" = ")
	out.WriteString(as.Value.String())
	out.WriteString(";")
	return out.String()
}
