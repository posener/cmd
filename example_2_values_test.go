package cmd_test

import (
	"fmt"

	"github.com/posener/complete/v2/predict"
	"github.com/posener/cmd"
)

// Flags and positional arguments can be defined with valid values. It is also possible to enable
// the check for valid on parsing times. When setting valid values, they will be completed by the
// bash completion system.
func Example_values() {
	// Should be defined in global `var`.
	var (
		root = cmd.New()
		// Define a flag with valid values 'foo' and 'bar', and enforce the values by `OptCheck()`.
		flag1 = root.String("flag1", "", "first flag", predict.OptValues("foo", "bar"), predict.OptCheck())
		// Define a flag with valid values of Go file names.
		file = root.String("file", "", "file path", predict.OptPredictor(predict.Files("*")), predict.OptCheck())
		// Define positional arguments with valid values 'baz' and 'buzz', and choose not to enforce
		// the check by not calling `OptCheck`.
		args = root.Args("[args...]", "positional arguments", predict.OptValues("baz", "buzz"))
	)

	// Should be in `main()`.
	root.Parse([]string{"cmd", "-flag1", "foo", "-file", "cmd.go", "buz", "bazz"})

	// Test:

	fmt.Println(*flag1, *file, *args)
	// Output: foo cmd.go [buz bazz]
}
