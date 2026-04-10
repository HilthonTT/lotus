package main

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	kindMethod   = 2
	kindFunction = 3
	kindVariable = 6
	kindModule   = 9
	kindKeyword  = 14
)

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}

type Analyzer struct {
	keywords []string
	builtins []builtinDoc
	packages map[string][]packageMember
}

type builtinDoc struct {
	name      string
	signature string
	doc       string
}

type packageMember struct {
	name      string
	signature string
	doc       string
}

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		keywords: []string{
			"let", "mut", "fn", "class", "extends", "if", "else",
			"while", "for", "in", "return", "break", "continue",
			"true", "false", "nil", "self", "super",
			"import", "export", "from",
		},
		builtins: []builtinDoc{
			{"print", "print(...values)", "Prints values to stdout separated by spaces."},
			{"len", "len(value) -> int", "Returns the length of a string, array, or map."},
			{"push", "push(array, value) -> array", "Returns a new array with value appended."},
			{"pop", "pop(array) -> value", "Returns the last element of the array."},
			{"head", "head(array) -> value", "Returns the first element of the array."},
			{"tail", "tail(array) -> array", "Returns all elements except the first."},
			{"type", "type(value) -> string", "Returns the type of a value as a string."},
			{"str", "str(value) -> string", "Converts a value to its string representation."},
			{"int", "int(value) -> int", "Converts a value to an integer."},
			{"range", "range([start,] end [,step]) -> array", "Returns an array of integers."},
		},
		packages: map[string][]packageMember{
			"Console": {
				{"readLine", "Console.readLine() -> string", "Reads a line from stdin."},
				{"prompt", "Console.prompt(message) -> string", "Prints message then reads a line from stdin."},
				{"clear", "Console.clear()", "Clears the terminal screen."},
				{"print", "Console.print(...values)", "Prints values to stdout."},
				{"printErr", "Console.printErr(...values)", "Prints values to stderr."},
			},
			"Math": {
				{"sqrt", "Math.sqrt(x) -> float", "Returns the square root of x."},
				{"abs", "Math.abs(x) -> number", "Returns the absolute value of x."},
				{"floor", "Math.floor(x) -> int", "Returns the floor of x as an integer."},
				{"pow", "Math.pow(base, exp) -> float", "Returns base raised to exp."},
				{"max", "Math.max(a, b) -> number", "Returns the larger of two values."},
				{"min", "Math.min(a, b) -> number", "Returns the smaller of two values."},
				{"pi", "Math.pi() -> float", "Returns π (3.141592653589793)."},
			},
			"OS": {
				{"exit", "OS.exit([code])", "Exits the process."},
				{"args", "OS.args() -> array", "Returns command-line arguments."},
				{"env", "OS.env(key) -> string | nil", "Returns an environment variable."},
				{"readFile", "OS.readFile(path) -> string | nil", "Reads a file as a string."},
				{"writeFile", "OS.writeFile(path, content) -> bool", "Writes a string to a file."},
				{"parseInt", "OS.parseInt(s) -> int | nil", "Parses a string to an integer."},
				{"parseFloat", "OS.parseFloat(s) -> float | nil", "Parses a string to a float."},
			},
		},
	}
}

func (a *Analyzer) Complete(source, prefix, receiver string) []CompletionItem {
	var items []CompletionItem

	if receiver != "" {
		if members, ok := a.packages[receiver]; ok {
			for _, m := range members {
				m := m
				items = append(items, CompletionItem{
					Label:         m.name,
					Kind:          intPtr(kindMethod),
					Detail:        strPtr(m.signature),
					Documentation: &MarkupContent{Kind: "markdown", Value: m.doc},
				})
			}
		}
		return items
	}

	for _, kw := range a.keywords {
		if strings.HasPrefix(kw, prefix) {
			kw := kw
			items = append(items, CompletionItem{Label: kw, Kind: intPtr(kindKeyword)})
		}
	}
	for _, b := range a.builtins {
		if strings.HasPrefix(b.name, prefix) {
			b := b
			items = append(items, CompletionItem{
				Label:         b.name,
				Kind:          intPtr(kindFunction),
				Detail:        strPtr(b.signature),
				Documentation: &MarkupContent{Kind: "markdown", Value: b.doc},
			})
		}
	}
	for name := range a.packages {
		if strings.HasPrefix(name, prefix) {
			name := name
			items = append(items, CompletionItem{Label: name, Kind: intPtr(kindModule)})
		}
	}
	for _, name := range extractUserSymbols(source) {
		if strings.HasPrefix(name, prefix) && name != prefix {
			name := name
			items = append(items, CompletionItem{Label: name, Kind: intPtr(kindVariable)})
		}
	}
	return items
}

func (a *Analyzer) HoverDoc(word string) string {
	for _, b := range a.builtins {
		if b.name == word {
			return fmt.Sprintf("```\n%s\n```\n\n%s", b.signature, b.doc)
		}
	}
	if _, ok := a.packages[word]; ok {
		return fmt.Sprintf("**%s** — built-in package", word)
	}
	return ""
}

var (
	reLetMut = regexp.MustCompile(`\b(?:let|mut)\s+([a-zA-Z_][a-zA-Z0-9_]*)`)
	reFn     = regexp.MustCompile(`\bfn\s+([a-zA-Z_][a-zA-Z0-9_]*)`)
	reClass  = regexp.MustCompile(`\bclass\s+([A-Z][a-zA-Z0-9_]*)`)
)

func extractUserSymbols(source string) []string {
	seen := map[string]bool{}
	var result []string
	add := func(name string) {
		if !seen[name] {
			seen[name] = true
			result = append(result, name)
		}
	}
	for _, m := range reLetMut.FindAllStringSubmatch(source, -1) {
		add(m[1])
	}
	for _, m := range reFn.FindAllStringSubmatch(source, -1) {
		add(m[1])
	}
	for _, m := range reClass.FindAllStringSubmatch(source, -1) {
		add(m[1])
	}
	return result
}
