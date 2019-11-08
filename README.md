# subcmd

[![Build Status](https://travis-ci.org/posener/subcmd.svg?branch=master)](https://travis-ci.org/posener/subcmd)
[![codecov](https://codecov.io/gh/posener/subcmd/branch/master/graph/badge.svg)](https://codecov.io/gh/posener/subcmd)
[![GoDoc](https://godoc.org/github.com/posener/subcmd?status.svg)](http://godoc.org/github.com/posener/subcmd)
[![goreadme](https://goreadme.herokuapp.com/badge/posener/subcmd.svg)](https://goreadme.herokuapp.com)

subcmd is a minimalistic library that enables easy sub commands with the standard `flag` library.

Define a `root` command object using the `Root` function.
This object exposes the standard library's `flag.FlagSet` API, which enables adding flags in the
standard way.
Additionally, this object exposes the `SubCommand` method, which returns another command object.
This objects also exposing the same API, enabling definition of flags and nested sub commands.

The root object then have to be called with the `Parse` or `ParseArgs` methods, similarly to
the `flag.Parse` call.

The usage is automatically configured to show both sub commands and flags.

#### Positional arguments

The `subcmd` library is opinionated about positional arguments: it enforces their definition
and parsing. The user can define for each sub command if and how many positional arguments it
accepts. Their usage is similar to the flag values usage.

#### Limitations

Suppose `cmd` has a flag `-flag`, and a subcommand `sub`. In the current implementation:
Calling `cmd sub -flag` won't work as the flag is set after the sub command, while
`cmd -flag sub` will work perfectly fine. Each flag needs to be used in the scope of its command.

#### Examples

Definition and usage of sub commands and sub commands flags.

```golang
package main

import (
	"fmt"

	"github.com/posener/subcmd"
)

var (
	// Define a root command. Some options can be set using the `Opt*` functions. It returns a
	// `*Cmd` object.
	root = subcmd.Root()
	// The `*Cmd` object can be used as the standard library `flag.FlagSet`.
	flag0 = root.String("flag0", "", "root stringflag")
	// From each command object, a sub command can be created. This can be done recursively.
	sub1 = root.SubCommand("sub1", "first sub command")
	// Each sub command can have flags attached.
	flag1 = sub1.String("flag1", "", "sub1 string flag")
	sub2  = root.SubCommand("sub2", "second sub command")
	flag2 = sub1.Int("flag2", 0, "sub2 int flag")
)

// Definition and usage of sub commands and sub commands flags.
func main() {
	// In the example we use `Parse()` for a given list of command line arguments. This is useful
	// for testing, but should be replaced with `root.ParseArgs()` in `main()`
	root.Parse([]string{"cmd", "sub1", "-flag1", "value"})

	// Usually the program should switch over the sub commands. The chosen sub command will return
	// true for the `Parsed()` method.
	switch {
	case sub1.Parsed():
		fmt.Printf("Called sub1 with flag: %s", *flag1)
	case sub2.Parsed():
		fmt.Printf("Called sub2 with flag: %d", *flag2)
	}
}

```


---

Created by [goreadme](https://github.com/apps/goreadme)
