// subcmd is a minimalistic library that enables easy sub commands with the standard `flag` library.
//
// Define a `root` command object using the `Root` function.
// This object exposes the standard library's `flag.FlagSet` API, which enables adding flags in the
// standard way.
// Additionally, this object exposes the `SubCommand` method, which returns another command object.
// This objects also exposing the same API, enabling definition of flags and nested sub commands.
//
// The root object then have to be called with the `Parse` or `ParseArgs` methods, similiraly to
// the `flag.Parse` call.
//
// The usage is automatically configured to show both sub commands and flags.
//
// Example
//
// See ./example/main.go.
//
// Limitations
//
// Suppose `cmd` has a flag `-flag`, and a subcommand `sub`. In the current implementation:
// Calling `cmd sub -flag` won't work as the flag is set after the sub command, while
// `cmd -flag sub` will work perfectly fine. Each flag needs to be used in the scope of its command.
package subcmd

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
)

// Cmd is a command that can have set of flags and a sub command.
type Cmd struct {
	name        string
	description string

	*flag.FlagSet
	subCmd map[string]*Cmd
}

type config struct {
	errorHandling flag.ErrorHandling
	output        io.Writer
	name          string
	description   string
}

// Option can configure command behavior.
type Option func(o *config)

// OptErrorHandling defines the behavior in case of an error in the `Parse` function.
func OptErrorHandling(errorHandling flag.ErrorHandling) Option {
	return func(cfg *config) {
		cfg.errorHandling = errorHandling
	}
}

// OptOutput sets the output for the usage.
func OptOutput(w io.Writer) Option {
	return func(cfg *config) {
		cfg.output = w
	}
}

// OptName sets a predefined name to the root command.
func OptName(name string) Option {
	return func(cfg *config) {
		cfg.name = name
	}
}

// OptDescription sets a description to the root command.
func OptDescription(description string) Option {
	return func(cfg *config) {
		cfg.description = description
	}
}

// Root creats a new root command.
func Root(options ...Option) *Cmd {
	// Set default config.
	cfg := config{
		errorHandling: flag.ExitOnError,
		output:        os.Stderr,
		name:          os.Args[0],
	}
	// Update with requested options.
	for _, option := range options {
		option(&cfg)
	}

	return newCmd(cfg)
}

// SubCommand creates a new sub command to the given command.
func (c *Cmd) SubCommand(name string, description string) *Cmd {
	if c.subCmd[name] != nil {
		panic(fmt.Sprintf("sub command %q already exists", name))
	}

	subCmd := newCmd(config{
		name:          name,
		description:   description,
		errorHandling: c.ErrorHandling(),
		output:        c.Output(),
	})
	c.subCmd[name] = subCmd
	return subCmd
}

// Parse command line arguments.
func (c *Cmd) ParseArgs() error {
	return c.Parse(os.Args)
}

// Parse a set of arguments.
func (c *Cmd) Parse(args []string) error {
	return c.handleError(c.parse(args))
}

func (c *Cmd) parse(args []string) error {
	if len(args) < 1 {
		panic("must be at least the command in arguments")
	}

	// Check for command flags, and update the remaining arguments.
	err := c.FlagSet.Parse(args[1:])
	if err != nil {
		return err
	}
	args = c.FlagSet.Args()

	// Check if another the first remaining argument matches any sub command.
	if len(args) == 0 {
		return nil
	}
	subCmd := c.subCmd[args[0]]
	if subCmd == nil {
		return nil
	}
	return subCmd.Parse(args)
}

func (c *Cmd) usage() {
	w := c.Output()
	fmt.Fprintf(w, "%s\t%s\n", c.name, c.description)
	if len(c.subCmd) > 0 {
		fmt.Fprintf(w, "Subcommands:\n")
		for _, name := range c.subCmdNames() {
			fmt.Fprintf(w, "  %s\t%s\n", name, c.subCmd[name].description)
		}
	}

	hasFlags := false
	c.FlagSet.VisitAll(func(*flag.Flag) { hasFlags = true })

	if hasFlags {
		fmt.Fprintf(w, "Flags:\n")
		c.FlagSet.PrintDefaults()
	}
}

// subCmdNames return all sub commands oredered alphabetically.
func (c *Cmd) subCmdNames() []string {
	subCmds := make([]string, 0, len(c.subCmd))
	for subCmd := range c.subCmd {
		subCmds = append(subCmds, subCmd)
	}
	sort.Strings(subCmds)
	return subCmds
}

func (c *Cmd) handleError(err error) error {
	if err == nil {
		return nil
	}
	switch c.ErrorHandling() {
	case flag.ExitOnError:
		os.Exit(2)
	case flag.PanicOnError:
		panic(err)
	}
	return err
}

func newCmd(cfg config) *Cmd {
	flagSet := flag.NewFlagSet(os.Args[0], cfg.errorHandling)
	flagSet.SetOutput(cfg.output)

	cmd := &Cmd{
		name:        cfg.name,
		description: cfg.description,
		FlagSet:     flagSet,
		subCmd:      make(map[string]*Cmd),
	}
	cmd.Usage = cmd.usage
	return cmd
}
