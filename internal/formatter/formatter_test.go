package formatter

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	loremIpsum          = "\nLorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.\n"
	formattedLoremIpsum = `  
  Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
  labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris
  nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit
  esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt
  in culpa qui officia deserunt mollit anim id est laborum.
`
)

func TestFormatter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		indent []byte
		width  int
		input  string
		want   string
	}{
		{
			name:   "indent",
			indent: []byte(" "),
			input:  "a\nbb\nccc",
			want:   " a\n bb\n ccc",
		},
		{
			name:  "width",
			width: 3,
			input: "a a bb",
			want:  "a a\nbb",
		},
		{
			name:   "indent and width",
			indent: []byte(" "),
			width:  3,
			input:  "aa b ccc",
			want:   " aa\n b\n ccc",
		},
		{
			name:  "line greater than width",
			width: 3,
			input: "aa bbbb cc",
			want:  "aa\nbbbb\ncc",
		},
		{
			name:  "width multiple spaces",
			width: 5,
			input: "  aa bbbb  cc  dddd",
			want:  "  aa\nbbbb \ncc \ndddd",
		},
		{
			name:   "lorem ipsum",
			indent: []byte("  "),
			width:  100,
			input:  loremIpsum,
			want:   formattedLoremIpsum,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got bytes.Buffer

			f := Formatter{Writer: &got, Indent: tt.indent, Width: tt.width}
			f.Write([]byte(tt.input))

			assert.Equal(t, tt.want, got.String())
		})
	}
}
