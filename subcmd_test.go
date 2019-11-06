package subcmd

import (
	"bytes"
	"errors"
	"flag"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testCommand struct {
	*Cmd
	rootFlag *bool

	sub1     *Cmd
	sub1Flag *string
	sub1Args *[]string

	sub11     *Cmd
	sub11Flag *string

	sub2     *Cmd
	sub2Args ArgsStr

	out bytes.Buffer
}

const longText = "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum."

func testRoot() *testCommand {
	var root testCommand

	root.Cmd = Root(
		OptName("cmd"),
		OptErrorHandling(flag.ContinueOnError),
		OptOutput(&root.out),
		OptSynopsis("cmd synopsis"),
		OptDetails("testing command line example"))

	root.rootFlag = root.Bool("flag", false, "example of bool flag")

	root.sub1 = root.SubCommand("sub1", "a sub command with flags and sub commands", OptDetails(longText))
	root.sub1Flag = root.sub1.String("flag", "", "example of string flag")
	root.sub1Args = root.sub1.Args("", "")

	root.sub11 = root.sub1.SubCommand("sub1", "sub command of sub command")
	root.sub11Flag = root.sub11.String("flag", "", "example of string flag")

	root.sub2 = root.SubCommand("sub2", "a sub command without flags and sub commands")
	root.sub2Args = make(ArgsStr, 0, 1)
	root.sub2.ArgsVar(&root.sub2Args, "[arg]", "arg is a single argument")

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
		sub1Args    []string
		sub11Flag   string
		sub2Args    []string
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
			args:    []string{"cmd", "arg1"},
			wantErr: true,
		},
		{
			args:       []string{"cmd", "sub1", "arg1"},
			sub1Parsed: true,
			sub1Args:   []string{"arg1"},
		},
		{
			args:       []string{"cmd", "sub1", "arg1", "arg2"},
			sub1Parsed: true,
			sub1Args:   []string{"arg1", "arg2"},
		},
		{
			args:       []string{"cmd", "sub1", "-flag", "value", "arg1"},
			sub1Parsed: true,
			sub1Flag:   "value",
			sub1Args:   []string{"arg1"},
		},
		{
			args:       []string{"cmd", "sub2", "arg1"},
			sub2Parsed: true,
		},
		{
			args:    []string{"cmd", "sub2", "arg1", "arg2"},
			wantErr: true,
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
			want: `Usage: cmd [flags]

cmd synopsis

  testing command line example

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
			want: `Usage: cmd sub1 [flags] [args...]

a sub command with flags and sub commands

  Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor
  incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
  nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.
  Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore
  eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt
  in culpa qui officia deserunt mollit anim id est laborum.

Subcommands:

  sub1	sub command of sub command

Flags:

  -flag string
    	example of string flag

`,
		},
		{
			args: []string{"cmd", "sub2", "-h"},
			want: `Usage: cmd sub2 [arg]

a sub command without flags and sub commands

Positional arguments:

  arg is a single argument

`,
		},
		{
			args: []string{"cmd", "sub1", "sub1", "-h"},
			want: `Usage: cmd sub1 sub1 [flags]

sub command of sub command

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
			assert.True(t, errors.As(err, &flag.ErrHelp))
			assert.Equal(t, tt.want, root.out.String())
		})
	}
}

func TestCmd_failures(t *testing.T) {
	t.Parallel()

	t.Run("command can't have two sub commands with the same name", func(t *testing.T) {
		cmd := Root()
		cmd.SubCommand("sub", "synopsis")

		assert.Panics(t, func() { cmd.SubCommand("sub", "synopsis") })
	})

	t.Run("parse must get at least one argument", func(t *testing.T) {
		cmd := Root()

		assert.Panics(t, func() { cmd.Parse(nil) })
	})

	t.Run("both command and sub command have positional arguments should panic", func(t *testing.T) {
		cmd := Root()
		cmd.Args("", "")
		subcmd := cmd.SubCommand("sub", "synopsis")
		subcmd.Args("", "")

		assert.Panics(t, func() { cmd.ParseArgs() })
	})

	t.Run("both command and sub sub command have positional arguments should panic", func(t *testing.T) {
		cmd := Root()
		cmd.Args("", "")
		sub := cmd.SubCommand("sub", "synopsis")
		subsub := sub.SubCommand("sub", "synopsis")
		subsub.Args("", "")

		assert.Panics(t, func() { cmd.ParseArgs() })
	})

	t.Run("both sub command and sub sub command have positional arguments should panic", func(t *testing.T) {
		cmd := Root()
		sub := cmd.SubCommand("sub", "synopsis")
		sub.Args("", "")
		subsub := sub.SubCommand("sub", "synopsis")
		subsub.Args("", "")

		assert.Panics(t, func() { cmd.ParseArgs() })
	})

	t.Run("two different sub command may have positional arguments", func(t *testing.T) {
		cmd := Root()
		sub1 := cmd.SubCommand("sub1", "synopsis")
		sub1.Args("", "")
		sub2 := cmd.SubCommand("sub2", "synopsis")
		sub2.Args("", "")

		assert.NotPanics(t, func() { cmd.Parse([]string{"cmd"}) })
	})

	t.Run("calling positional more than once is not allowed", func(t *testing.T) {
		cmd := Root()
		cmd.Args("", "")

		assert.Panics(t, func() { cmd.Args("", "") })
	})
}
