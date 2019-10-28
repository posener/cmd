package main

import (
	"log"

	"github.com/posener/subcmd"
)

// Define the commands and flags in the var section, similarly to the classic usage of the standard
// flag library.
var (
	// Define a root command. Some options can be set using the `Opt*` functions.
	root = subcmd.Root()
	// Define a flag for the root command, just as done in the standard flag library for a `FlagSet`
	// object.
	path = root.String("path", "", "path to write")

	// From each command object, a sub command can be created using the `SubCommand` method, giving
	// a name and a description.
	write = root.SubCommand("write", "write the file")
	// This sub command is exactly the same as any other command object, and can be used to define
	// sub command flags, or even nested sub commands.
	writeText = write.String("text", "", "text to write to file")

	// Define another sub command from the root command.
	read = root.SubCommand("read", "read the file")
)

func main() {
	// In the main function, call the `Parse` or `ParseArgs`, just as in the stasndard library
	// the `flag.Parse` function should be called.
	root.ParseArgs()

	// The flag variables behaves the same as in the stasndard library.
	if *path == "" {
		log.Fatal("Path must be provided")
	}

	// In order to understand which sub command was used, sub commands should be checked for their
	// `Parsed` method:
	switch {
	case write.Parsed():
		log.Printf("Writing %q to path %s", *writeText, *path)
		return
	case read.Parsed():
		log.Printf("Reading path %s", *path)
		return
	default:
		log.Fatal("no command was specidied.")
	}
}
