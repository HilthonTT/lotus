package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/hilthontt/lotus/compiler"
	"github.com/hilthontt/lotus/lexer"
	"github.com/hilthontt/lotus/object"
	"github.com/hilthontt/lotus/parser"
	"github.com/hilthontt/lotus/vm"
)

const Logo = `
██╗      ██████╗ ████████╗██╗   ██╗███████╗
██║     ██╔═══██╗╚══██╔══╝██║   ██║██╔════╝
██║     ██║   ██║   ██║   ██║   ██║███████╗
██║     ██║   ██║   ██║   ██║   ██║╚════██║
███████╗╚██████╔╝   ██║   ╚██████╔╝███████║
╚══════╝ ╚═════╝    ╚═╝    ╚═════╝ ╚══════╝

  A compiled language with a stack-based VM
  Type 'help' for commands, Ctrl+C to exit
`

const Oops = `
 ██████╗  ██████╗ ██████╗ ███████╗
██╔═══██╗██╔═══██╗██╔══██╗██╔════╝
██║   ██║██║   ██║██████╔╝███████╗
██║   ██║██║   ██║██╔═══╝ ╚════██║
╚██████╔╝╚██████╔╝██║     ███████║
 ╚═════╝  ╚═════╝ ╚═╝     ╚══════╝`

const prompt = "lotus » "
const contPrompt = "      … "

// GlobalsSize mirrors vm.GlobalsSize — keep in sync.
const GlobalsSize = 65536

// REPLState holds everything that must persist across inputs.
type REPLState struct {
	symbolTable *compiler.SymbolTable
	constants   []object.Object
	globals     []object.Object
}

func newREPLState() *REPLState {
	// Use a throwaway compiler to get the initial seeded symbol table
	// (all builtins + packages pre-defined at correct indices).
	seed := compiler.New()
	globals := make([]object.Object, GlobalsSize)

	// Seed built-in packages into the globals size
	for _, name := range compiler.BuiltinPackageOrder {
		pkg := compiler.BuiltinPackages[name]
		sym, ok := seed.PublicResolve(name)
		if ok {
			globals[sym.Index] = pkg
		}
	}

	return &REPLState{
		symbolTable: seed.ExportSymbolTable(),
		constants:   []object.Object{},
		globals:     globals,
	}
}

// Start runs the interactive REPL loop, reading from in and writing to out.
func Start(in io.Reader, out io.Writer, engine *string) {
	fmt.Fprint(out, Logo)
	fmt.Fprintf(out, "  Lotus REPL  (engine: %s)\n", *engine)
	fmt.Fprintln(out, "  Type 'exit' or Ctrl+D to quit.")

	scanner := bufio.NewScanner(in)
	state := newREPLState()

	for {
		fmt.Fprint(out, prompt)

		input, done := readInput(scanner, out)
		if done {
			fmt.Fprintln(out, "\nGoodbye! 🪷")
			return
		}
		if input == "" {
			continue
		}
		if strings.TrimSpace(input) == "exit" {
			fmt.Fprintln(out, "Goodbye! 🪷")
			return
		}

		if *engine == "vm" {
			evalVM(input, out, state)
		} else {
			evalEval(input, out)
		}
	}
}

// readInput reads one logical input from the user.
// If the user opens a brace block it keeps reading continuation lines
// until the braces are balanced.
func readInput(scanner *bufio.Scanner, out io.Writer) (input string, eof bool) {
	if !scanner.Scan() {
		return "", true
	}
	line := scanner.Text()
	if strings.TrimSpace(line) == "" {
		return "", false
	}

	var sb strings.Builder
	sb.WriteString(line)

	depth := braceDepth(line)
	for depth > 0 {
		fmt.Fprint(out, contPrompt)
		if !scanner.Scan() {
			break
		}
		next := scanner.Text()
		sb.WriteByte('\n')
		sb.WriteString(next)
		depth += braceDepth(next)
	}
	return sb.String(), false
}

// braceDepth counts the net brace nesting of a line, ignoring strings.
func braceDepth(line string) int {
	depth := 0
	inStr := false
	for _, ch := range line {
		if ch == '"' {
			inStr = !inStr
		}
		if inStr {
			continue
		}
		if ch == '{' {
			depth++
		}
		if ch == '}' {
			depth--
		}
	}
	return depth
}

// evalVM compiles and executes one REPL input using the bytecode VM,
// reusing and updating persistent state.
func evalVM(input string, out io.Writer, state *REPLState) {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	if errs := p.Errors(); len(errs) > 0 {
		for _, e := range errs {
			fmt.Fprintf(out, "  parse error: %s\n", e)
		}
		return
	}

	// Compile reusing the persistent symbol table and constants.
	comp := compiler.NewWithState(state.symbolTable, state.constants)
	if err := comp.Compile(program); err != nil {
		fmt.Fprintf(out, "  compile error: %s\n", err)
		return
	}

	// Save updated symbol table and constants for the next input.
	state.symbolTable = comp.ExportSymbolTable()
	state.constants = comp.Bytecode().Constants

	// Run using the persistent globals slice.
	machine := vm.NewWithGlobalsState(comp.Bytecode(), state.globals)
	if err := machine.Run(); err != nil {
		fmt.Fprintf(out, "  error: %s\n", err)
		return
	}

	// Print the last expression's value if it's non-nil.
	result := machine.LastPoppedStackElement()
	if result != nil && result.Type() != object.NIL_OBJ {
		fmt.Fprintf(out, "  => %s\n", result.Inspect())
	}
}

// evalEval runs one input through the tree-walking evaluator.
// State is not persisted across inputs in eval mode.
func evalEval(_ string, out io.Writer) {
	fmt.Fprintln(out, "  (persistent state requires --engine vm)")
}
