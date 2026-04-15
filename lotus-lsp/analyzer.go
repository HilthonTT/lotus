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
	kindField    = 5
	kindClass    = 7
)

type Analyzer struct {
	keywords  []string
	typeNames []string
	builtins  []builtinDoc
	packages  map[string][]packageMember
}

type builtinDoc struct {
	name, signature, doc string
}

type packageMember struct {
	name, signature, doc string
}

// classInfo holds extracted fields and methods for a class.
type classInfo struct {
	fields  []string
	methods []string
}

var (
	reLetMut    = regexp.MustCompile(`\b(?:let|mut)\s+([a-zA-Z_][a-zA-Z0-9_]*)`)
	reFn        = regexp.MustCompile(`\bfn\s+([a-zA-Z_][a-zA-Z0-9_]*)`)
	reClass     = regexp.MustCompile(`\bclass\s+([A-Z][a-zA-Z0-9_]*)`)
	reVarClass  = regexp.MustCompile(`\b(?:let|mut)\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*(?::\s*[a-zA-Z_][a-zA-Z0-9_]*)?\s*=\s*([A-Z][a-zA-Z0-9_]*)\s*\(`)
	reSelfField = regexp.MustCompile(`\bself\.([a-zA-Z_][a-zA-Z0-9_]*)\s*=`)
	reMethod    = regexp.MustCompile(`\bfn\s+([a-zA-Z_][a-zA-Z0-9_]*)\s*\(\s*self`)
)

func NewAnalyzer() *Analyzer {
	return &Analyzer{
		keywords: []string{
			"let", "mut", "fn", "class", "extends", "if", "else",
			"while", "for", "in", "return", "break", "continue",
			"true", "false", "nil", "self", "super",
			"import", "export", "from", "match", "enum",
		},
		typeNames: []string{
			"int", "float", "string", "bool", "array", "map", "nil",
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
				{"ceil", "Math.ceil(x) -> int", "Returns the ceiling of x as an integer."},
				{"round", "Math.round(x) -> int", "Rounds x to the nearest integer."},
				{"pow", "Math.pow(base, exp) -> float", "Returns base raised to exp."},
				{"max", "Math.max(a, b) -> number", "Returns the larger of two values."},
				{"min", "Math.min(a, b) -> number", "Returns the smaller of two values."},
				{"pi", "Math.pi() -> float", "Returns π (3.141592653589793)."},
				{"e", "Math.e() -> float", "Returns e (2.718281828459045)."},
			},
			"OS": {
				{"exit", "OS.exit([code])", "Exits the process with an optional exit code."},
				{"args", "OS.args() -> array", "Returns command-line arguments as an array."},
				{"env", "OS.env(key) -> string | nil", "Returns an environment variable by key."},
				{"readFile", "OS.readFile(path) -> string | nil", "Reads a file and returns its contents as a string."},
				{"writeFile", "OS.writeFile(path, content) -> bool", "Writes a string to a file. Returns true on success."},
				{"parseInt", "OS.parseInt(s) -> int | nil", "Parses a string to an integer."},
				{"parseFloat", "OS.parseFloat(s) -> float | nil", "Parses a string to a float."},
			},
			"Task": {
				{"spawn", "Task.spawn(fn())", "Runs a zero-argument Lotus closure in a new goroutine."},
				{"spawnWith", "Task.spawnWith(fn(arg), arg)", "Runs a Lotus closure in a new goroutine, passing one argument."},
				{"wait", "Task.wait()", "Blocks until all spawned tasks have finished."},
				{"sleep", "Task.sleep(ms: int)", "Pauses the current task for the given number of milliseconds."},
				{"mutex", "Task.mutex() -> Mutex", "Creates and returns a new mutex object."},
			},
			"String": {
				{"split", "String.split(str, sep) -> array", "Splits a string by separator into an array."},
				{"trim", "String.trim(str) -> string", "Removes leading and trailing whitespace."},
				{"trimLeft", "String.trimLeft(str) -> string", "Removes leading whitespace."},
				{"trimRight", "String.trimRight(str) -> string", "Removes trailing whitespace."},
				{"upper", "String.upper(str) -> string", "Converts string to uppercase."},
				{"lower", "String.lower(str) -> string", "Converts string to lowercase."},
				{"replace", "String.replace(str, old, new) -> string", "Replaces all occurrences of old with new."},
				{"contains", "String.contains(str, substr) -> bool", "Returns true if str contains substr."},
				{"startsWith", "String.startsWith(str, prefix) -> bool", "Returns true if str starts with prefix."},
				{"endsWith", "String.endsWith(str, suffix) -> bool", "Returns true if str ends with suffix."},
				{"indexOf", "String.indexOf(str, substr) -> int", "Returns the index of substr in str, or -1."},
				{"repeat", "String.repeat(str, n) -> string", "Repeats str n times."},
				{"padLeft", "String.padLeft(str, n, char) -> string", "Pads str on the left to length n."},
				{"padRight", "String.padRight(str, n, char) -> string", "Pads str on the right to length n."},
				{"chars", "String.chars(str) -> array", "Returns an array of single-character strings."},
				{"len", "String.len(str) -> int", "Returns the character count of str."},
				{"join", "String.join(array, sep) -> string", "Joins array elements with separator."},
				{"slice", "String.slice(str, start, end) -> string", "Returns a substring from start to end."},
			},
			"Array": {
				{"filter", "Array.filter(arr, fn(elem) -> bool) -> array", "Returns elements for which fn returns true."},
				{"map", "Array.map(arr, fn(elem) -> any) -> array", "Transforms each element using fn."},
				{"reduce", "Array.reduce(arr, fn(acc, elem) -> any, initial) -> any", "Reduces array to a single value."},
				{"find", "Array.find(arr, fn(elem) -> bool) -> elem | nil", "Returns the first matching element."},
				{"findIndex", "Array.findIndex(arr, fn(elem) -> bool) -> int", "Returns the index of the first match, or -1."},
				{"forEach", "Array.forEach(arr, fn(elem))", "Calls fn for each element."},
				{"contains", "Array.contains(arr, value) -> bool", "Returns true if array contains value."},
				{"reverse", "Array.reverse(arr) -> array", "Returns a reversed copy of the array."},
				{"sort", "Array.sort(arr) -> array", "Returns a sorted copy of the array."},
				{"sortBy", "Array.sortBy(arr, fn(elem) -> comparable) -> array", "Sorts by a key function."},
				{"flat", "Array.flat(arr) -> array", "Flattens one level of nested arrays."},
				{"join", "Array.join(arr, sep) -> string", "Joins array elements into a string."},
				{"slice", "Array.slice(arr, start, end) -> array", "Returns a sub-array."},
				{"unique", "Array.unique(arr) -> array", "Removes duplicate elements."},
				{"len", "Array.len(arr) -> int", "Returns the length of the array."},
				{"any", "Array.any(arr, fn(elem) -> bool) -> bool", "Returns true if any element matches."},
				{"all", "Array.all(arr, fn(elem) -> bool) -> bool", "Returns true if all elements match."},
			},
			"Time": {
				{"now", "Time.now() -> int", "Returns current time as Unix milliseconds."},
				{"sleep", "Time.sleep(ms: int)", "Pauses execution for ms milliseconds."},
				{"format", "Time.format(ms: int, layout: string) -> string", "Formats a timestamp. Layout: \"2006-01-02 15:04:05\""},
				{"parse", "Time.parse(str: string, layout: string) -> int", "Parses a time string into Unix milliseconds."},
				{"since", "Time.since(ms: int) -> int", "Milliseconds elapsed since the given timestamp."},
				{"until", "Time.until(ms: int) -> int", "Milliseconds until the given timestamp."},
				{"add", "Time.add(ms: int, duration: int) -> int", "Adds duration milliseconds to a timestamp."},
				{"diff", "Time.diff(a: int, b: int) -> int", "Returns (a - b) in milliseconds."},
				{"year", "Time.year(ms: int) -> int", "Extracts the year from a timestamp."},
				{"month", "Time.month(ms: int) -> int", "Extracts the month (1-12) from a timestamp."},
				{"day", "Time.day(ms: int) -> int", "Extracts the day (1-31) from a timestamp."},
				{"hour", "Time.hour(ms: int) -> int", "Extracts the hour (0-23) from a timestamp."},
				{"minute", "Time.minute(ms: int) -> int", "Extracts the minute (0-59) from a timestamp."},
				{"second", "Time.second(ms: int) -> int", "Extracts the second (0-59) from a timestamp."},
				{"weekday", "Time.weekday(ms: int) -> string", "Returns the weekday name e.g. \"Monday\"."},
				{"unix", "Time.unix(ms: int) -> int", "Converts Unix milliseconds to Unix seconds."},
				{"fromUnix", "Time.fromUnix(sec: int) -> int", "Converts Unix seconds to Unix milliseconds."},
				{"isBefore", "Time.isBefore(a: int, b: int) -> bool", "Returns true if timestamp a is before b."},
				{"isAfter", "Time.isAfter(a: int, b: int) -> bool", "Returns true if timestamp a is after b."},
				{"startOfDay", "Time.startOfDay(ms: int) -> int", "Returns the midnight timestamp for the given day."},
				{"endOfDay", "Time.endOfDay(ms: int) -> int", "Returns the 23:59:59.999 timestamp for the given day."},
				{"addDays", "Time.addDays(ms: int, days: int) -> int", "Adds n calendar days to a timestamp."},
				{"addMonths", "Time.addMonths(ms: int, months: int) -> int", "Adds n calendar months to a timestamp."},
				{"addYears", "Time.addYears(ms: int, years: int) -> int", "Adds n calendar years to a timestamp."},
				{"duration", "Time.duration(ms: int) -> string", "Human-readable duration e.g. \"2h 35m 10s\"."},
				{"ms", "Time.ms(n: int) -> int", "Returns n as milliseconds."},
				{"seconds", "Time.seconds(n: int) -> int", "Returns n seconds expressed in milliseconds."},
				{"minutes", "Time.minutes(n: int) -> int", "Returns n minutes expressed in milliseconds."},
				{"hours", "Time.hours(n: int) -> int", "Returns n hours expressed in milliseconds."},
				{"days", "Time.days(n: int) -> int", "Returns n days expressed in milliseconds."},
				{"utc", "Time.utc(ms: int) -> int", "Converts local timestamp to UTC."},
				{"timezone", "Time.timezone() -> string", "Returns the local timezone name."},
			},
			"Json": {
				{"stringify", "Json.stringify(value) -> string", "Converts a Lotus value to a JSON string."},
				{"prettyPrint", "Json.prettyPrint(value) -> string", "Converts a Lotus value to indented JSON."},
				{"parse", "Json.parse(str: string) -> value", "Parses a JSON string into a Lotus value (objects→Hash, arrays→Array)."},
				{"valid", "Json.valid(str: string) -> bool", "Returns true if the string is valid JSON."},
				{"keys", "Json.keys(str: string) -> array", "Returns the top-level keys of a JSON object string."},
				{"get", "Json.get(str: string, key: string) -> value", "Gets a top-level key from a JSON object string."},
				{"set", "Json.set(str: string, key: string, value) -> string", "Sets a key in a JSON object, returns new JSON string."},
				{"merge", "Json.merge(a: string, b: string) -> string", "Merges two JSON objects (b overwrites a on conflicts)."},
			},
		},
	}
}

// extractClasses parses class bodies from source and returns a map of
// className -> classInfo (fields set via self.x = ..., and methods).
func extractClasses(source string) map[string]classInfo {
	classes := map[string]classInfo{}

	// Find each class block: class Foo { ... }
	reClassBlock := regexp.MustCompile(`(?s)\bclass\s+([A-Z][a-zA-Z0-9_]*)[^{]*\{(.*?)\n\}`)
	for _, m := range reClassBlock.FindAllStringSubmatch(source, -1) {
		name := m[1]
		body := m[2]

		seen := map[string]bool{}
		info := classInfo{}

		for _, fm := range reSelfField.FindAllStringSubmatch(body, -1) {
			field := fm[1]
			if !seen[field] {
				seen[field] = true
				info.fields = append(info.fields, field)
			}
		}
		for _, mm := range reMethod.FindAllStringSubmatch(body, -1) {
			method := mm[1]
			if method != "init" {
				info.methods = append(info.methods, method)
			}
		}
		classes[name] = info
	}
	return classes
}

// extractVarTypes returns a map of varName -> className for lines like:
// let v = Vector(3.0, 4.0)
func extractVarTypes(source string) map[string]string {
	varTypes := map[string]string{}
	for _, m := range reVarClass.FindAllStringSubmatch(source, -1) {
		varTypes[m[1]] = m[2]
	}
	return varTypes
}

func (a *Analyzer) Complete(source, prefix, receiver string) []CompletionItem {
	var items []CompletionItem

	if receiver != "" {
		// Built-in packages
		if members, ok := a.packages[receiver]; ok {
			for _, m := range members {
				m := m
				items = append(items, CompletionItem{
					Label:         m.name,
					Kind:          new(kindMethod),
					Detail:        new(m.signature),
					Documentation: &MarkupContent{Kind: "markdown", Value: m.doc},
				})
			}
			return items
		}

		// Instance fields and methods
		varTypes := extractVarTypes(source)
		classes := extractClasses(source)
		className, ok := varTypes[receiver]
		if !ok {
			className = receiver
		}
		if info, ok := classes[className]; ok {
			for _, field := range info.fields {
				items = append(items, CompletionItem{
					Label:  field,
					Kind:   new(kindField),
					Detail: new(className + "." + field),
				})
			}
			for _, method := range info.methods {
				items = append(items, CompletionItem{
					Label:  method,
					Kind:   new(kindMethod),
					Detail: new(className + "." + method + "(self, ...)"),
				})
			}
			return items
		}

		return items
	}

	// Keywords
	for _, kw := range a.keywords {
		if strings.HasPrefix(kw, prefix) {
			kw := kw
			items = append(items, CompletionItem{Label: kw, Kind: new(kindKeyword)})
		}
	}

	// Type names (shown as keywords in completion)
	for _, t := range a.typeNames {
		if strings.HasPrefix(t, prefix) {
			t := t
			items = append(items, CompletionItem{
				Label:         t,
				Kind:          new(kindKeyword),
				Detail:        new("type"),
				Documentation: &MarkupContent{Kind: "markdown", Value: fmt.Sprintf("Built-in type: `%s`", t)},
			})
		}
	}

	// Builtins
	for _, b := range a.builtins {
		if strings.HasPrefix(b.name, prefix) {
			b := b
			items = append(items, CompletionItem{
				Label:         b.name,
				Kind:          new(kindFunction),
				Detail:        new(b.signature),
				Documentation: &MarkupContent{Kind: "markdown", Value: b.doc},
			})
		}
	}

	// Packages
	for name := range a.packages {
		if strings.HasPrefix(name, prefix) {
			name := name
			items = append(items, CompletionItem{Label: name, Kind: new(kindModule)})
		}
	}

	// User symbols
	for _, name := range extractUserSymbols(source) {
		if strings.HasPrefix(name, prefix) && name != prefix {
			name := name
			items = append(items, CompletionItem{Label: name, Kind: new(kindVariable)})
		}
	}

	return items
}

func (a *Analyzer) HoverDoc(word string) string {
	// Builtins
	for _, b := range a.builtins {
		if b.name == word {
			return fmt.Sprintf("```\n%s\n```\n\n%s", b.signature, b.doc)
		}
	}

	// Packages
	if members, ok := a.packages[word]; ok {
		var sb strings.Builder
		fmt.Fprintf(&sb, "**%s** — built-in package\n\n", word)
		sb.WriteString("| Member | Signature |\n|--------|----------|\n")
		for _, m := range members {
			fmt.Fprintf(&sb, "| `%s` | `%s` |\n", m.name, m.signature)
		}
		return sb.String()
	}

	// Type names
	typeDoc := map[string]string{
		"int":    "Built-in integer type. Example: `let x: int = 42`",
		"float":  "Built-in float type. Example: `let x: float = 3.14`",
		"string": "Built-in string type. Example: `let s: string = \"hello\"`",
		"bool":   "Built-in boolean type. Values: `true` or `false`",
		"array":  "Built-in array type. Example: `let a: array = [1, 2, 3]`",
		"map":    "Built-in map type. Example: `let m: map = {\"key\": \"value\"}`",
	}
	if doc, ok := typeDoc[word]; ok {
		return doc
	}

	return ""
}

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
