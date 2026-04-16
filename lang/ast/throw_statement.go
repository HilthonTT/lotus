package ast

import "github.com/hilthontt/lotus/token"

// ThrowStatement: throw "something went wrong"
type ThrowStatement struct {
	Token token.Token
	Value Expression
}

func (ts *ThrowStatement) statementNode()       {}
func (ts *ThrowStatement) TokenLiteral() string { return ts.Token.Literal }
func (ts *ThrowStatement) String() string {
	return "throw " + ts.Value.String()
}
