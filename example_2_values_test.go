package subcmd_test

import (
	"fmt"

	"github.com/posener/complete/v2/predict"
	"github.com/posener/subcmd"
)

// Flags and positional arguments can be defined with valid values. It is also possible to enable
// the check for valid on parsing times. When setting valid values, they will be completed by the
// bash completion system.
func Example_values() {
	// Should be defined in global `var`.
	var (
		cmd = subcmd.New()
		// Define a flag with valid values 'foo' and 'bar', and enforce the values by `OptCheck()`.
		flag1 = cmd.String("flag1", "", "first flag", predict.OptValues("foo", "bar"), predict.OptCheck())
		// Define a flag with valid values of Go file names.
		file = cmd.String("file", "", "file path", predict.OptPredictor(predict.Files("*")), predict.OptCheck())
		// Define positional arguments with valid values 'baz' and 'buzz', and choose not to enforce
		// the check by not calling `OptCheck`.
		args = cmd.Args("[args...]", "positional arguments", predict.OptValues("baz", "buzz"))
	)

	// Should be in `main()`.
	cmd.Parse([]string{"cmd", "-flag1", "foo", "-file", "subcmd.go", "buz", "bazz"})

	// Test:

	fmt.Println(*flag1, *file, *args)
	// Output: foo subcmd.go [buz bazz]
}
