package ast

type TypeParam struct {
	Name       string
	Constraint string // optional interface name, "" if unconstrained
}
