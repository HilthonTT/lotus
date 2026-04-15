package compiler

import (
	"math/rand/v2"

	"github.com/hilthontt/lotus/object"
)

func mathPackage() *object.Package {
	return &object.Package{
		Name: "Math",
		Functions: map[string]object.PackageFunction{
			"sqrt": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				v := toFloat64(args[0])
				return &object.Float{Value: mathSqrt(v)}
			},
			"abs": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				switch a := args[0].(type) {
				case *object.Integer:
					v := a.Value
					if v < 0 {
						v = -v
					}
					return &object.Integer{Value: v}
				case *object.Float:
					v := a.Value
					if v < 0 {
						v = -v
					}
					return &object.Float{Value: v}
				}
				return &object.Nil{}
			},
			"floor": func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return &object.Nil{}
				}
				return &object.Integer{Value: int64(toFloat64(args[0]))}
			},
			"pow": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}

				base := toFloat64(args[0])
				exp := toFloat64(args[1])
				return &object.Float{Value: mathPow(base, exp)}
			},
			"max": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}
				a, b := toFloat64(args[0]), toFloat64(args[1])
				if a > b {
					return args[0] // Return a
				}
				return args[1] // Return b
			},
			"min": func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return &object.Nil{}
				}
				a, b := toFloat64(args[0]), toFloat64(args[1])
				if a < b {
					return args[0] // Return b
				}
				return args[1] // Return a
			},
			"pi": func(args ...object.Object) object.Object {
				return &object.Float{Value: 3.141592653589793}
			},
			"random": func(args ...object.Object) object.Object {
				randomNum := rand.Float64()
				return &object.Float{Value: randomNum}
			},
		},
	}
}

func mathSqrt(x float64) float64 {
	if x < 0 {
		return 0
	}
	z := x / 2
	for range 100 {
		z -= (z*z - x) / (2 * z)
	}
	return z
}

func mathPow(base, exp float64) float64 {
	result := 1.0
	for exp > 0 {
		if int(exp)%2 == 0 {
			result *= base
		}
		base *= base
		exp = float64(int(exp) / 2)
	}
	return result
}

func toFloat64(o object.Object) float64 {
	switch v := o.(type) {
	case *object.Integer:
		return float64(v.Value)
	case *object.Float:
		return v.Value
	}
	return 0
}
