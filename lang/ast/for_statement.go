package ast

import (
	"bytes"

	"github.com/hilthontt/lotus/token"
)

type ForStatement struct {
	Token    token.Token
	Variable *Identifier
	Iterable Expression
	Body     *BlockStatement
}

func (s *ForStatement) statementNode() {}

// TokenLiteral returns the literal token.
func (s *ForStatement) TokenLiteral() string {
	return s.Token.Literal
}

func (s *ForStatement) String() string {
	var out bytes.Buffer
	out.WriteString("for ")
	out.WriteString(s.Variable.String())
	out.WriteString(" in ")
	out.WriteString(s.Iterable.String())
	out.WriteString(" { ")
	out.WriteString(s.Body.String())
	out.WriteString(" }")
	return out.String()
}
