package compiler

import (
	"encoding/json"
	"strings"

	"github.com/hilthontt/lotus/object"
)

func jsonPackage() *object.Package {
	return &object.Package{
		Name: "Json",
		Functions: map[string]object.PackageFunction{

			// Json.stringify(value) -> string
			// Converts a Lotus value to a JSON string.
			"stringify": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				native := toNative(args[0])
				bytes, err := json.Marshal(native)
				if err != nil {
					return &object.Nil{}
				}
				return &object.String{Value: string(bytes)}
			},

			// Json.prettyPrint(value) -> string
			// Converts a Lotus value to an indented JSON string.
			"prettyPrint": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				native := toNative(args[0])
				bytes, err := json.MarshalIndent(native, "", "  ")
				if err != nil {
					return &object.Nil{}
				}
				return &object.String{Value: string(bytes)}
			},

			// Json.parse(str: string) -> value
			// Parses a JSON string into a Lotus value.
			// JSON objects  → Hash
			// JSON arrays   → Array
			// JSON strings  → String
			// JSON numbers  → Integer or Float
			// JSON booleans → Boolean
			// JSON null     → nil
			"parse": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				s, ok := args[0].(*object.String)
				if !ok {
					return &object.Nil{}
				}
				var raw any
				if err := json.Unmarshal([]byte(s.Value), &raw); err != nil {
					return &object.Nil{}
				}
				return fromNative(raw)
			},

			// Json.valid(str: string) -> bool
			// Returns true if the string is valid JSON.
			"valid": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Boolean{Value: false}
				}
				s, ok := args[0].(*object.String)
				if !ok {
					return &object.Boolean{Value: false}
				}
				return &object.Boolean{Value: json.Valid([]byte(s.Value))}
			},

			// Json.keys(str: string) -> array
			// Returns the top-level keys of a JSON object string.
			"keys": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				s, ok := args[0].(*object.String)
				if !ok {
					return &object.Nil{}
				}
				var m map[string]any
				if err := json.Unmarshal([]byte(s.Value), &m); err != nil {
					return &object.Nil{}
				}
				elems := make([]object.Object, 0, len(m))
				for k := range m {
					elems = append(elems, &object.String{Value: k})
				}
				return &object.Array{Elements: elems}
			},

			// Json.get(str: string, key: string) -> value
			// Gets a top-level key from a JSON object string.
			"get": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}
				s, ok1 := args[0].(*object.String)
				k, ok2 := args[1].(*object.String)
				if !ok1 || !ok2 {
					return &object.Nil{}
				}
				var m map[string]any
				if err := json.Unmarshal([]byte(s.Value), &m); err != nil {
					return &object.Nil{}
				}
				val, exists := m[k.Value]
				if !exists {
					return &object.Nil{}
				}
				return fromNative(val)
			},

			// Json.set(str: string, key: string, value) -> string
			// Sets a top-level key in a JSON object and returns the new JSON string.
			"set": func(args ...object.Object) object.Object {
				if len(args) != 3 {
					return &object.Nil{}
				}
				s, ok1 := args[0].(*object.String)
				k, ok2 := args[1].(*object.String)
				if !ok1 || !ok2 {
					return &object.Nil{}
				}
				var m map[string]any
				if err := json.Unmarshal([]byte(s.Value), &m); err != nil {
					m = make(map[string]any)
				}
				m[k.Value] = toNative(args[2])
				bytes, err := json.Marshal(m)
				if err != nil {
					return &object.Nil{}
				}
				return &object.String{Value: string(bytes)}
			},

			// Json.merge(a: string, b: string) -> string
			// Merges two JSON objects (b overwrites a on conflicts).
			"merge": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}
				sa, ok1 := args[0].(*object.String)
				sb, ok2 := args[1].(*object.String)
				if !ok1 || !ok2 {
					return &object.Nil{}
				}
				var ma, mb map[string]any
				if err := json.Unmarshal([]byte(sa.Value), &ma); err != nil {
					ma = make(map[string]any)
				}
				if err := json.Unmarshal([]byte(sb.Value), &mb); err != nil {
					mb = make(map[string]any)
				}
				for k, v := range mb {
					ma[k] = v
				}
				bytes, err := json.Marshal(ma)
				if err != nil {
					return &object.Nil{}
				}
				return &object.String{Value: string(bytes)}
			},
		},
	}
}

// toNative converts a Lotus object to a Go native value for JSON marshalling.
func toNative(obj object.Object) any {
	switch o := obj.(type) {
	case *object.Integer:
		return o.Value
	case *object.Float:
		return o.Value
	case *object.Boolean:
		return o.Value
	case *object.String:
		return o.Value
	case *object.Nil:
		return nil
	case *object.Array:
		result := make([]any, len(o.Elements))
		for i, el := range o.Elements {
			result[i] = toNative(el)
		}
		return result
	case *object.Hash:
		result := make(map[string]any)
		for _, pair := range o.Pairs {
			key := pair.Key.Inspect()
			// Strip quotes if the key is a string literal
			key = strings.Trim(key, "\"")
			result[key] = toNative(pair.Value)
		}
		return result
	case *object.Instance:
		result := make(map[string]any)
		for k, v := range o.Fields {
			result[k] = toNative(v)
		}
		return result
	case *object.EnumVariant:
		if len(o.Data) == 0 {
			return o.EnumName + "." + o.VariantName
		}
		result := map[string]any{
			"type": o.EnumName + "." + o.VariantName,
		}
		for k, v := range o.Data {
			result[k] = toNative(v)
		}
		return result
	default:
		return obj.Inspect()
	}
}

// fromNative converts a Go native value (from JSON unmarshalling) to a Lotus object.
func fromNative(v any) object.Object {
	if v == nil {
		return &object.Nil{}
	}
	switch val := v.(type) {
	case bool:
		return &object.Boolean{Value: val}
	case float64:
		// JSON numbers are float64 — use integer if it's a whole number
		if val == float64(int64(val)) {
			return &object.Integer{Value: int64(val)}
		}
		return &object.Float{Value: val}
	case string:
		return &object.String{Value: val}
	case []any:
		elems := make([]object.Object, len(val))
		for i, el := range val {
			elems[i] = fromNative(el)
		}
		return &object.Array{Elements: elems}
	case map[string]any:
		pairs := make(map[object.HashKey]object.HashPair)
		for k, v := range val {
			key := &object.String{Value: k}
			pairs[key.HashKey()] = object.HashPair{
				Key:   key,
				Value: fromNative(v),
			}
		}
		return &object.Hash{Pairs: pairs}
	default:
		return &object.String{Value: "null"}
	}
}
