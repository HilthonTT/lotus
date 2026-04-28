package compiler

import (
	"encoding/gob"
	"fmt"
	"os"

	"github.com/hilthontt/lotus/code"
	"github.com/hilthontt/lotus/object"
)

func init() {
	// Register every concrete type that may appear in the constants pool
	// so gob can encode and decode them correctly.
	gob.Register(&object.Integer{})
	gob.Register(&object.Float{})
	gob.Register(&object.String{})
	gob.Register(&object.Boolean{})
	gob.Register(&object.Nil{})
	gob.Register(&object.Array{})
	gob.Register(&object.Hash{})
	gob.Register(object.HashPair{})
	gob.Register(object.HashKey{})
	gob.Register(&object.CompiledFunction{})
	gob.Register(&object.EnumDef{})
	gob.Register(&object.EnumVariantDef{})
	gob.Register(&object.Interface{})
	gob.Register(object.InterfaceMethodSpec{})
	gob.Register(code.LineTable{})
	gob.Register(code.LineEntry{})
}

const (
	bcMagic   uint32 = 0x4C4F5442 // "LOTB"
	bcVersion uint8  = 2          // bump when format changes incompatibly
)

type serializedBytecode struct {
	Magic           uint32
	FormatVersion   uint8
	Instructions    code.Instructions
	Constants       []object.Object
	ExportedSymbols map[string]int
	Lines           code.LineTable
}

// WriteBytecode serializes bc to a .lotusbc file at path.
func WriteBytecode(bc *Bytecode, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create %q: %w", path, err)
	}
	defer f.Close()

	sb := serializedBytecode{
		Magic:           bcMagic,
		FormatVersion:   bcVersion,
		Instructions:    []byte(bc.Instructions),
		Constants:       bc.Constants,
		ExportedSymbols: bc.ExportedSymbols,
		Lines:           bc.Lines,
	}
	if err := gob.NewEncoder(f).Encode(sb); err != nil {
		return fmt.Errorf("encode bytecode: %w", err)
	}
	return nil
}

// ReadBytecode deserializes a .lotusbc file into a Bytecode.
func ReadBytecode(path string) (*Bytecode, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open %q: %w", path, err)
	}
	defer f.Close()

	var sb serializedBytecode
	if err := gob.NewDecoder(f).Decode(&sb); err != nil {
		return nil, fmt.Errorf("decode %q: %w", path, err)
	}
	if sb.Magic != bcMagic {
		return nil, fmt.Errorf("%q is not a valid .lotusbc file", path)
	}
	if sb.FormatVersion != bcVersion {
		return nil, fmt.Errorf(
			"bytecode version mismatch in %q: file=%d runtime=%d — recompile required",
			path, sb.FormatVersion, bcVersion,
		)
	}

	return &Bytecode{
		Instructions:    code.Instructions(sb.Instructions),
		Constants:       sb.Constants,
		ExportedSymbols: sb.ExportedSymbols,
		Lines:           sb.Lines,
	}, nil
}
