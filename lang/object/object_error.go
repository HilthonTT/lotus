package object

const ERROR_OBJ ObjectType = "ERROR"

// LotusError is a thrown error value — distinct from a Go runtime error.
// It's what gets bound to the catch variable.
type LotusError struct {
	Message string
}

func (e *LotusError) Type() ObjectType {
	return ERROR_OBJ
}

func (e *LotusError) Inspect() string {
	return e.Message
}
