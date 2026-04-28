package code

import "encoding/gob"

func init() {
	gob.Register(LineTable{})
	gob.Register(LineEntry{})
}

// LineEntry maps a contiguous range of bytecode bytes to a source line.
type LineEntry struct {
	Start int // first byf offset of this run
	Count int // number of bytes on this line
	Line  int // 1-based source line number
}

// LineTable is a compact run-length-encoded offset -> line mapping
type LineTable []LineEntry

// LineForOffset returns the 1-based source line for a bytecode offset.
// Returns 0 if the offset is not covered.
func (lt LineTable) LineForOffset(offset int) int {
	for _, e := range lt {
		if offset >= e.Start && offset < e.Start+e.Count {
			return e.Line
		}
	}
	return 0
}
