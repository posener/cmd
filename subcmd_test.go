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
	sub2Args *[]string

	out bytes.Buffer
}

func testRoot() *testCommand {
	var root testCommand

	root.Cmd = Root(
		OptName("cmd"),
		OptErrorHandling(flag.ContinueOnError),
		OptOutput(&root.out),
		OptSynopsis("cmd synopsys"),
		OptDetails("testing command line example"))

	root.rootFlag = root.Bool("flag", false, "example of bool flag")

	root.sub1 = root.SubCommand("sub1", "a sub command with flags and sub commands", OptDetails("sub command details"))
	root.sub1Flag = root.sub1.String("flag", "", "example of string flag")
	root.sub1Args = root.sub1.Args(nil)

	root.sub11 = root.sub1.SubCommand("sub1", "sub command of sub command")
	root.sub11Flag = root.sub11.String("flag", "", "example of string flag")

	root.sub2 = root.SubCommand("sub2", "a sub command without flags and sub commands")
	root.sub2Args = root.sub2.Args(&ArgsOpts{N: 1, Usage: "[arg]", Details: "arg is a single argument"})

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

cmd synopsys

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
			want: `Usage: cmd sub1 [flags] [args]

a sub command with flags and sub commands

  sub command details

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
		cmd.SubCommand("sub", "synopsys")

		assert.Panics(t, func() { cmd.SubCommand("sub", "synopsys") })
	})

	t.Run("parse must get at least one argument", func(t *testing.T) {
		cmd := Root()

		assert.Panics(t, func() { cmd.Parse(nil) })
	})

	t.Run("both command and sub command have positional arguments should panic", func(t *testing.T) {
		cmd := Root()
		cmd.Args(nil)
		subcmd := cmd.SubCommand("sub", "synopsys")
		subcmd.Args(nil)

		assert.Panics(t, func() { cmd.ParseArgs() })
	})

	t.Run("both command and sub sub command have positional arguments should panic", func(t *testing.T) {
		cmd := Root()
		cmd.Args(nil)
		sub := cmd.SubCommand("sub", "synopsys")
		subsub := sub.SubCommand("sub", "synopsys")
		subsub.Args(nil)

		assert.Panics(t, func() { cmd.ParseArgs() })
	})

	t.Run("both sub command and sub sub command have positional arguments should panic", func(t *testing.T) {
		cmd := Root()
		sub := cmd.SubCommand("sub", "synopsys")
		sub.Args(nil)
		subsub := sub.SubCommand("sub", "synopsys")
		subsub.Args(nil)

		assert.Panics(t, func() { cmd.ParseArgs() })
	})

	t.Run("two different sub command may have positional arguments", func(t *testing.T) {
		cmd := Root()
		sub1 := cmd.SubCommand("sub1", "synopsys")
		sub1.Args(nil)
		sub2 := cmd.SubCommand("sub2", "synopsys")
		sub2.Args(nil)

		assert.NotPanics(t, func() { cmd.Parse([]string{"cmd"}) })
	})

	t.Run("calling positional more than once is not allowed", func(t *testing.T) {
		cmd := Root()
		cmd.Args(nil)

		assert.Panics(t, func() { cmd.Args(nil) })
	})
}
