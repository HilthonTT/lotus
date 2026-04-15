package compiler

import (
	"os"
	"strconv"

	"github.com/hilthontt/lotus/object"
)

func osPackage() *object.Package {
	return &object.Package{
		Name: "OS",
		Functions: map[string]object.PackageFunction{
			"exit": func(args ...object.Object) object.Object {
				code := 0
				if len(args) == 1 {
					if i, ok := args[0].(*object.Integer); ok {
						code = int(i.Value)
					}
				}
				os.Exit(code)
				return &object.Nil{}
			},
			"args": func(args ...object.Object) object.Object {
				elems := make([]object.Object, len(os.Args))
				for i, a := range os.Args {
					elems[i] = &object.String{Value: a}
				}
				return &object.Array{Elements: elems}
			},
			"env": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				key := args[0].Inspect()
				val := os.Getenv(key)
				if val == "" {
					return &object.Nil{}
				}
				return &object.String{Value: val}
			},
			"readFile": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				data, err := os.ReadFile(args[0].Inspect())
				if err != nil {
					return &object.Nil{}
				}
				return &object.String{Value: string(data)}
			},
			"writeFile": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}
				err := os.WriteFile(args[0].Inspect(), []byte(args[1].Inspect()), 0644)
				if err != nil {
					return &object.Boolean{Value: false}
				}
				return &object.Boolean{Value: true}
			},
			"parseInt": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				i, err := strconv.ParseInt(args[0].Inspect(), 10, 64)
				if err != nil {
					return &object.Nil{}
				}
				return &object.Integer{Value: i}
			},
			"parseFloat": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				f, err := strconv.ParseFloat(args[0].Inspect(), 64)
				if err != nil {
					return &object.Nil{}
				}
				return &object.Float{Value: f}
			},
		},
	}
}
