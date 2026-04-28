package vm

import (
	"fmt"
	"path/filepath"
	"strings"
)

// StackFrame is one entry in a printed stack trace.
type StackFrame struct {
	FuncName string
	File     string
	Line     int
}

// BuildStackTrace walks the live frame stack from innermost to outermost
// and resolves each frame's current source line via its LineTable.
func (vm *VM) BuildStackTrace() []StackFrame {
	var frames []StackFrame

	for i := vm.framesIndex - 1; i >= 0; i-- {
		f := vm.frames[i]
		if f == nil {
			continue
		}

		fn := f.closure.Fn
		name := fn.Name
		if name == "" {
			if i == 0 {
				name = "<main>"
			} else {
				name = "<anonymous>"
			}
		}

		ip := f.ip
		if ip < 0 {
			ip = 0
		}

		line := 0
		if fn.Lines != nil {
			line = fn.Lines.LineForOffset(ip)
		}

		frames = append(frames, StackFrame{
			FuncName: name,
			File:     vm.filePath,
			Line:     line,
		})
	}

	return frames
}

// FormatStackTrace returns a human-readable stack trace string.
// Example:
//
//	Stack trace (most recent call first):
//	  at divide (main.lotus:12)
//	  at safeDivide (main.lotus:27)
//	  at <main> (main.lotus:45)
func FormatStackTrace(frames []StackFrame) string {
	if len(frames) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("Stack trace (most recent call first):\n")
	for _, f := range frames {
		base := filepath.Base(f.File)
		if f.Line > 0 {
			sb.WriteString(fmt.Sprintf("  at %s (%s:%d)\n", f.FuncName, base, f.Line))
		} else {
			sb.WriteString(fmt.Sprintf("  at %s (%s)\n", f.FuncName, base))
		}
	}
	return sb.String()
}

// runtimeError wraps err with the current stack trace.
// If no file path or no frames, returns err unchanged.
func (vm *VM) runtimeError(err error) error {
	if vm.filePath == "" {
		return err
	}
	frames := vm.BuildStackTrace()
	if len(frames) == 0 {
		return err
	}
	trace := FormatStackTrace(frames)
	return fmt.Errorf("%s\n%s", err.Error(), trace)
}

// SetFilePath stores the source file name used in stack trace output.
func (vm *VM) SetFilePath(path string) {
	vm.filePath = path
}
