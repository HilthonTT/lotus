package object

type ObjectType string

const (
	INTEGER_OBJ ObjectType = "INTEGER"
	FLOAT_OBJ   ObjectType = "FLOAT"
	BOOLEAN_OBJ ObjectType = "BOOLEAN"
	STRING_OBJ  ObjectType = "STRING"
	NIL_OBJ     ObjectType = "NIL"
	ARRAY_OBJ   ObjectType = "ARRAY"
	MAP_OBJ     ObjectType = "MAP"
	CLOSURE_OBJ ObjectType = "CLOSURE"
	BUILTIN_OBJ ObjectType = "BUILTIN"
)

// Object is the interface that all of our various object-types must implemented.
type Object interface {
	// Type returns the type of this object.
	Type() ObjectType

	// Inspect returns a string-representation of the given object.
	Inspect() string
}
