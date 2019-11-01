package formatter

import (
	"io"
)

// Formatter is an `io.Writer` that enables indenting text and wrapping it to a specified width.
type Formatter struct {
	// Writer to write the formatted text.
	io.Writer
	// Indent the text with a defined byte slice.
	Indent []byte
	// Width is the line size for wrapping the text.
	Width int

	curWidth  int
	lastSpace int
	hasSpace  bool
}

func (f *Formatter) Write(b []byte) (int, error) {
	f.validate()
	b = f.insertIndent(b)
	return f.Writer.Write(b)
}

func (f *Formatter) insertIndent(b []byte) []byte {
	for i := 0; i < len(b); i++ {
		// Insert indentation if a new line.
		if len(f.Indent) > 0 && f.curWidth == 0 {
			i, b = insert(b, i, f.Indent)
			f.curWidth = len(f.Indent)
		} else {
			f.curWidth++
		}

		switch b[i] {
		case '\n':
			f.curWidth = 0
			f.lastSpace = i
		case ' ', '\t':
			f.lastSpace = i
			f.hasSpace = true
		default:
			if f.Width > 0 && f.curWidth > f.Width {
				if f.hasSpace {
					b[f.lastSpace] = '\n'
					f.hasSpace = false
					i = f.lastSpace - 1 // start next loop from the new line.
				}
			}
		}
	}
	return b
}

func insert(buf []byte, i int, in []byte) (int, []byte) {
	before, after := buf[:i], buf[i:]
	buf = append(before, append(in, after...)...)
	return i + len(in), buf
}

func (f *Formatter) validate() {
	if len(f.Indent) > 0 && f.Width > 0 && f.Width <= len(f.Indent) {
		panic("Formatter: width must be greater than indent length.")
	}
}
