package ast

import "github.com/hilthontt/lotus/token"

// The Call is wrapped in a zero-arg closure at parse time so it executes lazily.
type DeferStatement struct {
	Token token.Token // the 'defer' token
	Call  Expression
}

func (ds *DeferStatement) statementNode() {}

func (ds *DeferStatement) TokenLiteral() string {
	return ds.Token.Literal
}

func (ds *DeferStatement) String() string {
	return "defer " + ds.Call.String()
}
