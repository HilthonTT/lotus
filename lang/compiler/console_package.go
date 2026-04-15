package compiler

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hilthontt/lotus/object"
)

func consolePackage() *object.Package {
	return &object.Package{
		Name: "Console",
		Functions: map[string]object.PackageFunction{
			// readLine() -> string  — reads a full line from stdin
			"readLine": func(args ...object.Object) object.Object {
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					return &object.String{Value: scanner.Text()}
				}
				return &object.String{Value: ""}
			},

			// readLine with prompt: Console.prompt("Name: ") -> string
			"prompt": func(args ...object.Object) object.Object {
				if len(args) == 1 {
					fmt.Print(args[0].Inspect())
				}
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					return &object.String{Value: scanner.Text()}
				}
				return &object.String{Value: ""}
			},

			// clear() — clears the terminal screen
			"clear": func(args ...object.Object) object.Object {
				fmt.Print("\033[H\033[2J")
				return &object.Nil{}
			},

			// print(...) — same as the global print but namespaced
			"print": func(args ...object.Object) object.Object {
				parts := make([]string, len(args))
				for i, a := range args {
					parts[i] = a.Inspect()
				}
				fmt.Println(strings.Join(parts, " "))
				return &object.Nil{}
			},

			// printErr(...) — writes to stderr
			"printErr": func(args ...object.Object) object.Object {
				parts := make([]string, len(args))
				for i, a := range args {
					parts[i] = a.Inspect()
				}
				fmt.Fprintln(os.Stderr, strings.Join(parts, " "))
				return &object.Nil{}
			},
		},
	}
}
