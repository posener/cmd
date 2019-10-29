package subcmd

import (
	"bytes"
	"flag"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCommand struct {
	*Cmd
	sub1      *Cmd
	sub11     *Cmd
	sub2      *Cmd
	rootFlag  *bool
	sub1Flag  *string
	sub11Flag *string
	out       bytes.Buffer
}

func testRoot() *testCommand {
	var root testCommand

	root.Cmd = Root(
		OptName("cmd"),
		OptDescription("description"),
		OptErrorHandling(flag.ContinueOnError),
		OptOutput(&root.out))

	root.rootFlag = root.Bool("flag", false, "example of bool flag")

	root.sub1 = root.SubCommand("sub1", "a sub command with flags and sub commands")
	root.sub1Flag = root.sub1.String("flag", "", "example of string flag")

	root.sub11 = root.sub1.SubCommand("sub1", "sub command of sub command")
	root.sub11Flag = root.sub11.String("flag", "", "example of string flag")

	root.sub2 = root.SubCommand("sub2", "a sub command without flags and sub commands")

	return &root
}

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
			args:       []string{"cmd", "sub2"},
			sub2Parsed: true,
		},
		{
			args:    []string{"cmd", "-no-such-flag"},
			wantErr: true,
		},
		{
			args:    []string{"cmd", "-rootflag", "-no-such-flag"},
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
			root := testRoot()
			err := root.Parse(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.sub1Parsed, root.sub1.Parsed())
				assert.Equal(t, tt.sub11Parsed, root.sub11.Parsed())
				assert.Equal(t, tt.sub2Parsed, root.sub2.Parsed())
				assert.Equal(t, tt.rootFlag, *root.rootFlag)
				assert.Equal(t, tt.sub1Flag, *root.sub1Flag)
				assert.Equal(t, tt.sub11Flag, *root.sub11Flag)
			}
		})
	}
}

func TestHelp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		args []string
		want string
	}{
		{
			args: []string{"cmd", "-h"},
			want: `cmd	description
Subcommands:
  sub1	a sub command with flags and sub commands
  sub2	a sub command without flags and sub commands
Flags:
  -flag
    	example of bool flag
`,
		},
		{
			args: []string{"cmd", "sub1", "-h"},
			want: `sub1	a sub command with flags and sub commands
Subcommands:
  sub1	sub command of sub command
Flags:
  -flag string
    	example of string flag
`,
		},
		{
			args: []string{"cmd", "sub2", "-h"},
			want: `sub2	a sub command without flags and sub commands
`,
		},
		{
			args: []string{"cmd", "sub1", "sub1", "-h"},
			want: `sub1	sub command of sub command
Flags:
  -flag string
    	example of string flag
`,
		},
	}

	for _, tt := range tests {
		t.Run(strings.Join(tt.args, " "), func(t *testing.T) {
			root := testRoot()
			err := root.Parse(tt.args)
			assert.Equal(t, err, flag.ErrHelp)
			assert.Equal(t, tt.want, root.out.String())
		})
	}
}
