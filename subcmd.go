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
// Positional arguments
//
// The `subcmd` library is opinionated about positional arguments: it enforces their definition
// and parsing. The user can define for each sub command if and how many positional arguments it
// accepts. Their usage is similar to the flag values usage.
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
	"strings"
)

// Cmd is a command that can have set of flags and a sub command.
type Cmd struct {
	// name of command or sub command.
	name string
	// synopsis of the command.
	synopsis string
	// argsOpts are the positional arguments options. If nil the command does not accept
	// positional arguments.
	argsOpts *ArgsOpts

	// FlagsSet holds the flags of the command.
	*flag.FlagSet
	// args holds the positional arguments of the commands.
	args *[]string
	// sub holds the sub commands of the command.
	sub map[string]*Cmd
}

// ArgsOpts are options for positional arguments.
type ArgsOpts struct {
	// N can be used to enforce a fixed number of positional arguments. Any non-positive number will
	// be ignored.
	N int
	// Usage is a string representing the positional arguments which will be printed in the command
	// usage. For example, The string "[source] [destination]" can represent usage of positional
	// arguments for a move command.
	Usage string
	// Details can be used to provide farther explaination on the positional arguments in the usage
	// help.
	Details string
}

type config struct {
	errorHandling flag.ErrorHandling
	output        io.Writer
	name          string
	synopsis      string
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

// OptSynopsis sets a description to the root command.
func OptSynopsis(synopsis string) Option {
	return func(cfg *config) {
		cfg.synopsis = synopsis
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

// Args returns the positional arguments for the command and enable defining options. Only a sub
// command that called this method accepts positional arguments. Calling a sub command with
// positional arguments where they were not defined result in parsing error. The provided options
// can be nil for default values.
func (c *Cmd) Args(opts *ArgsOpts) *[]string {
	if c.argsOpts != nil {
		panic("Args() called more than once.")
	}
	// Default options.
	if opts == nil {
		opts = &ArgsOpts{}
	}
	c.argsOpts = opts
	return c.args
}

// SubCommand creates a new sub command to the given command.
func (c *Cmd) SubCommand(name string, synopsis string) *Cmd {
	if c.sub[name] != nil {
		panic(fmt.Sprintf("sub command %q already exists", name))
	}

	subCmd := newCmd(config{
		name:          name,
		synopsis:      synopsis,
		errorHandling: c.ErrorHandling(),
		output:        c.Output(),
	})
	c.sub[name] = subCmd
	return subCmd
}

// Parse command line arguments.
func (c *Cmd) ParseArgs() error {
	return c.Parse(os.Args)
}

// Parse a set of arguments.
func (c *Cmd) Parse(args []string) error {
	c.validate()
	_, err := c.parse(args)

	return c.handleError(err)
}

func (c *Cmd) parse(args []string) ([]string, error) {
	if len(args) < 1 {
		panic("must be at least the command in arguments")
	}

	// Check for command flags, and update the remaining arguments.
	err := c.FlagSet.Parse(args[1:])
	if err != nil {
		return nil, fmt.Errorf("%s: bad flags: %w", c.name, err)
	}
	args = c.FlagSet.Args()

	// Check if another the first remaining argument matches any sub command.
	if len(args) > 0 && c.sub[args[0]] != nil {
		subcmd := c.sub[args[0]]
		args, err = subcmd.parse(args)
		if err != nil {
			return nil, fmt.Errorf("%s > %v", c.name, err)
		}
	}

	// Collect positional arguments if required.
	args, err = c.setArgs(args)
	if err != nil {
		return nil, fmt.Errorf("%s: bad positional args: %v", c.name, err)
	}

	return args, nil
}

func (c *Cmd) setArgs(args []string) ([]string, error) {
	opt := c.argsOpts
	if opt == nil {
		if len(args) > 0 {
			return nil, fmt.Errorf("positional args not expected, got %v", args)
		}
		return args, nil
	}
	if opt.N > 0 && len(args) != opt.N {
		return nil, fmt.Errorf("required %d positional args, got %v", opt.N, args)
	}
	c.args = &args
	return nil, nil
}

// validate the command line. Panics on error.
func (c *Cmd) validate() {
	c.validatePositional(false)
}

// validatePositional validates positional arguments. If c was defined with positional arguments,
// any of its sub commands can't be defined with positional arguments.
func (c *Cmd) validatePositional(foundPositional bool, chain ...string) {
	chain = append(chain, c.name)
	hasPositional := c.argsOpts != nil

	if foundPositional && hasPositional {
		panic("A command with positional arguments can't have a sub command with positional arguments. Found chain: " + strings.Join(chain, ">"))
	}

	for _, subcmd := range c.sub {
		subcmd.validatePositional(foundPositional || hasPositional, chain...)
	}
}

func (c *Cmd) usage() {
	w := c.Output()
	usage := "Usage: " + c.name
	if c.hasFlags() {
		usage += " [flags]"
	}
	if c.argsOpts != nil {
		usage += " " + c.argsOpts.usage()
	}
	fmt.Fprintln(w, usage)
	fmt.Fprintln(w, c.synopsis)
	if len(c.sub) > 0 {
		fmt.Fprintf(w, "Subcommands:\n")
		for _, name := range c.subNames() {
			fmt.Fprintf(w, "  %s\t%s\n", name, c.sub[name].synopsis)
		}
	}

	if c.hasFlags() {
		fmt.Fprintf(w, "Flags:\n")
		c.FlagSet.PrintDefaults()
	}

	if c.argsOpts != nil && c.argsOpts.Details != "" {
		fmt.Fprintf(w, "Positional arguments:\n%s\n", c.argsOpts.Details)
	}
}

// subNames return all sub commands oredered alphabetically.
func (c *Cmd) subNames() []string {
	names := make([]string, 0, len(c.sub))
	for name := range c.sub {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

func (c *Cmd) hasFlags() bool {
	hasFlags := false
	c.FlagSet.VisitAll(func(*flag.Flag) { hasFlags = true })
	return hasFlags
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

func (o *ArgsOpts) usage() string {
	if o.Usage != "" {
		return o.Usage
	}
	return "[args]"
}

func newCmd(cfg config) *Cmd {
	flagSet := flag.NewFlagSet(os.Args[0], cfg.errorHandling)
	flagSet.SetOutput(cfg.output)

	cmd := &Cmd{
		name:     cfg.name,
		synopsis: cfg.synopsis,
		FlagSet:  flagSet,
		sub:      make(map[string]*Cmd),
	}
	cmd.Usage = cmd.usage
	return cmd
}
