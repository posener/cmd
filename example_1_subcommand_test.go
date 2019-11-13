package subcmd_test

import (
	"fmt"

	"github.com/posener/subcmd"
)

var (
	// Define a cmd command. Some options can be set using the `Opt*` functions. It returns a
	// `*Cmd` object.
	cmd = subcmd.New()
	// The `*Cmd` object can be used as the standard library `flag.FlagSet`.
	flag0 = cmd.String("flag0", "", "root string flag")

	// From each command object, a sub command can be created. This can be done recursively.
	sub1 = cmd.SubCommand("sub1", "first sub command")
	// Each sub command can have flags attached.
	flag1 = sub1.String("flag1", "", "sub1 string flag")

	sub2  = cmd.SubCommand("sub2", "second sub command")
	flag2 = sub1.Int("flag2", 0, "sub2 int flag")
)

// Definition and usage of sub commands and sub commands flags.
func Example() {
	// In the example we use `Parse()` for a given list of command line arguments. This is useful
	// for testing, but should be replaced with `cmd.ParseArgs()` in `main()`
	cmd.Parse([]string{"cmd", "sub1", "-flag1", "value"})

	// Usually the program should switch over the sub commands. The chosen sub command will return
	// true for the `Parsed()` method.
	switch {
	case sub1.Parsed():
		fmt.Printf("Called sub1 with flag: %s", *flag1)
	case sub2.Parsed():
		fmt.Printf("Called sub2 with flag: %d", *flag2)
	}
	// Output: Called sub1 with flag: value
}
