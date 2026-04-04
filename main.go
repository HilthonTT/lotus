package main

import (
	"fmt"
	"os"

	"github.com/hilthontt/lotus/code"
	"github.com/hilthontt/lotus/compiler"
	"github.com/hilthontt/lotus/lexer"
	"github.com/hilthontt/lotus/parser"
	"github.com/hilthontt/lotus/version"
	"github.com/hilthontt/lotus/vm"
)

const logo = `
██╗      ██████╗ ████████╗██╗   ██╗███████╗
██║     ██╔═══██╗╚══██╔══╝██║   ██║██╔════╝
██║     ██║   ██║   ██║   ██║   ██║███████╗
██║     ██║   ██║   ██║   ██║   ██║╚════██║
███████╗╚██████╔╝   ██║   ╚██████╔╝███████║
╚══════╝ ╚═════╝    ╚═╝    ╚═════╝ ╚══════╝

  A compiled language with a stack-based VM
  Type 'help' for commands, Ctrl+C to exit
`

func main() {
	repl()

	runFile("./examples/example.lotus", false)
}

func repl() {
	v := version.GetVersionInfo()
	fmt.Print(logo, v.Version)
	fmt.Println()
}

func runFile(filename string, disassemble bool) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %s\n", err)
		os.Exit(1)
	}

	input := string(data)
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	errors := p.Errors()
	if len(errors) > 0 {
		for _, e := range errors {
			fmt.Fprintf(os.Stderr, "parse error: %s\n", e)
		}
		os.Exit(1)
	}

	c := compiler.New()
	if err := c.Compile(program); err != nil {
		fmt.Fprintf(os.Stderr, "compile error: %s\n", err)
		os.Exit(1)
	}

	bytecode := c.Bytecode()

	if disassemble {
		fmt.Println("=== Disassembly ===")
		fmt.Print(code.Disassemble(bytecode.Instructions))
		fmt.Println("\n=== Constants ===")
		for i, k := range bytecode.Constants {
			fmt.Printf("  [%d] %s (%s)\n", i, k.Inspect(), k.Type())
		}
		fmt.Println()
		return
	}

	machine := vm.New(bytecode)
	if err := machine.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %s\n", err)
		os.Exit(1)
	}
}
