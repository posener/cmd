package subcmd

import (
	"bytes"
	"flag"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubCmd(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args        []string
		wantErr     bool
		sub1Parsed  bool
		sub11Parsed bool
		sub2Parsed  bool
		rootFlag    bool
		sub1Flag    string
		sub11Flag   string
		sub2Flag    int
		wantArgs    []string
	}{
		{
			args: []string{"cmd"},
		},
		{
			args:     []string{"cmd", "-flag"},
			rootFlag: true,
		},
		{
			args:       []string{"cmd", "-flag", "sub1", "-flag", "value"},
			rootFlag:   true,
			sub1Parsed: true,
			sub1Flag:   "value",
		},
		{
			args:       []string{"cmd", "sub1", "-flag", "value"},
			sub1Parsed: true,
			sub1Flag:   "value",
		},
		{
			args:        []string{"cmd", "sub1", "sub1", "-flag", "value"},
			sub1Parsed:  true,
			sub11Parsed: true,
			sub11Flag:   "value",
		},
		{
			args:       []string{"cmd", "sub2", "-flag", "1"},
			sub2Parsed: true,
			sub2Flag:   1,
		},
		{
			args:    []string{"cmd", "-no-such-flag"},
			wantErr: true,
		},
		{
			args:    []string{"cmd", "-flag", "-no-such-flag"},
			wantErr: true,
		},
		{
			args:    []string{"cmd", "sub1", "-no-such-flag"},
			wantErr: true,
		},
		{
			args:     []string{"cmd", "arg1"},
			wantArgs: []string{"arg1"},
		},
		{
			args:       []string{"cmd", "sub1", "arg1"},
			sub1Parsed: true,
			wantArgs:   []string{"arg1"},
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {

			root := Root(OptErrorHandling(flag.ContinueOnError), OptOutput(ioutil.Discard))
			rootFlag := root.Bool("flag", false, "example of bool flag")

			sub1 := root.SubCommand("sub1", "sub command of root command")
			sub1Flag := sub1.String("flag", "", "example of string flag")

			sub11 := sub1.SubCommand("sub1", "sub command of sub command")
			sub11Flag := sub11.String("flag", "", "example of string flag")

			sub2 := root.SubCommand("sub2", "another sub command of root command")
			sub2Flag := sub2.Int("flag", 0, "example of int flag")

			err := root.Parse(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.sub1Parsed, sub1.Parsed())
				assert.Equal(t, tt.sub11Parsed, sub11.Parsed())
				assert.Equal(t, tt.sub2Parsed, sub2.Parsed())
				assert.Equal(t, tt.rootFlag, *rootFlag)
				assert.Equal(t, tt.sub1Flag, *sub1Flag)
				assert.Equal(t, tt.sub11Flag, *sub11Flag)
				assert.Equal(t, tt.sub2Flag, *sub2Flag)
			}
		})
	}
}

func TestHelp(t *testing.T) {
	t.Parallel()

	rootHelp := `cmd	description
Subcommands:
  sub1	sub command of root command
  sub2	another sub command of root command
`
	sub1Help := `sub1	sub command of root command
Subcommands:
  sub1	sub command of sub command
Flags:
  -flag string
    	example of string flag
`

	tests := []struct {
		args []string
		want string
	}{
		{
			args: []string{"cmd", "-h"},
			want: rootHelp,
		},
		{
			args: []string{"cmd", "sub1", "-h"},
			want: sub1Help,
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			var out bytes.Buffer

			root := Root(
				OptName("cmd"),
				OptDescription("description"),
				OptErrorHandling(flag.ContinueOnError),
				OptOutput(&out))

			sub1 := root.SubCommand("sub1", "sub command of root command")
			_ = sub1.String("flag", "", "example of string flag")

			sub11 := sub1.SubCommand("sub1", "sub command of sub command")
			_ = sub11.String("flag", "", "example of string flag")

			sub2 := root.SubCommand("sub2", "another sub command of root command")
			_ = sub2.Int("flag", 0, "example of int flag")

			err := root.Parse(tt.args)
			assert.Equal(t, err, flag.ErrHelp)
			assert.Equal(t, tt.want, out.String())
		})
	}
}
