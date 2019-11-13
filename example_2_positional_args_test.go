package subcmd_test

import (
	"fmt"

	"github.com/posener/subcmd"
)

// Usage of positional arguments. If a program accepts positional arguments it must declare it using
// the `Args()` or the `ArgsVar()` methods. Positional arguments can be also defined on sub
// commands.
func Example_args() {
	// Should be defined in global `var`.
	var (
		cmd = subcmd.New()
		// Positional arguments can be defined as any other flag.
		args = cmd.Args("[args...]", "positional arguments for command line")
	)

	// Should be in `main()`.
	cmd.Parse([]string{"cmd", "v1", "v2", "v3"})

	// Test:

	fmt.Println(*args)
	// Output: [v1 v2 v3]
}

// Usage of positional arguments with exact number of arguments.
func Example_argsN() {
	// Should be defined in global `var`.
	var (
		cmd = subcmd.New()
		// Define arguments with cap=2 will ensure that the number of arguments is always 2.
		args = make(subcmd.ArgsStr, 2)
	)

	// Should be in `init()`.
	cmd.ArgsVar(&args, "[src] [dst]", "positional arguments for command line")

	// Should be in `main()`.
	cmd.Parse([]string{"cmd", "from.txt", "to.txt"})

	// Test:

	fmt.Println(args)
	// Output: [from.txt to.txt]
}

// Usage of positional arguments of a specific type.
func Example_argsInt() {
	// Should be defined in global `var`.
	var (
		cmd = subcmd.New()
		// Define positional arguments of type integer.
		args subcmd.ArgsInt
	)

	// Should be in `init()`.
	cmd.ArgsVar(&args, "[int...]", "numbers to sum")

	// Should be in `main()`.
	cmd.Parse([]string{"cmd", "10", "20", "30"})

	// Test:

	sum := 0
	for _, n := range args {
		sum += n
	}
	fmt.Println(sum)
	// Output: 60
}

// Usage of positional arguments with a conversion function.
func Example_argsFn() {
	// Should be defined in global `var`.
	var (
		cmd      = subcmd.New()
		src, dst string
	)

	// A function that convert the positional arguments to the program variables.
	argsFn := func(args []string) error {
		if len(args) != 2 {
			return fmt.Errorf("expected src and dst, got %d arguments", len(args))
		}
		src, dst = args[0], args[1]
		return nil
	}

	// Should be in `init()`.
	cmd.ArgsVar(subcmd.ArgsFn(argsFn), "[src] [dst]", "positional arguments for command line")

	// Should be in `main()`.
	cmd.Parse([]string{"cmd", "from.txt", "to.txt"})

	// Test:

	fmt.Println(src, dst)
	// Output: from.txt to.txt
}
